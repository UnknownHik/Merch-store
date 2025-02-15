package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"API-Avito-shop/config"
	"API-Avito-shop/internal/delivery"
	"API-Avito-shop/internal/middleware"
	"API-Avito-shop/internal/repositories"
	"API-Avito-shop/internal/services"
	"API-Avito-shop/internal/utils/logger"
	_ "API-Avito-shop/internal/utils/validation"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	dbPool    *pgxpool.Pool
	config    *config.Config
	logger    *slog.Logger
	apiServer *http.Server
}

// New создает приложение
func New() (*App, error) {
	app := &App{}

	// Загружаем конфигурацию
	cfg, err := config.MustLoad()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	app.config = cfg

	// Инициализируем логгер
	app.logger = logger.InitLogger(slog.LevelDebug)

	// Подключаемся к базе данных
	pool, err := newDBConn(&app.config.DatabaseConfig)
	if err != nil {
		app.logger.Error("Failed to connect to the database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	app.dbPool = pool

	// Настройка API-сервера
	if err = app.setupAPIServer(); err != nil {
		return nil, fmt.Errorf("failed to setup API server: %w", err)
	}

	app.logger.Info("Application initialized successfully")
	return app, nil
}

// Run запускает приложение
func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		app.logger.Info("API server started successfully", "address", app.apiServer.Addr)
		if err := app.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Failed to start the server", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	app.logger.Info("Received shutdown signal")

	return app.shutdown()
}

// shutdown останавливает сервер и закрывает соединения
func (app *App) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if app.apiServer != nil {
		if err := app.apiServer.Shutdown(ctx); err != nil {
			app.logger.Error("HTTP server shutdown failed", "error", err)
			return err
		}
	}

	if app.dbPool != nil {
		app.dbPool.Close()
		app.logger.Info("Database connection closed successfully")
	}

	app.logger.Info("Application stopped gracefully")
	return nil
}

// setupAPIServer настраивает HTTP-сервер
func (app *App) setupAPIServer() error {
	secretKey := app.config.ApiServerConfig.AuthSecretKey
	token := services.NewToken(secretKey, app.logger)

	// Инициализация репозитория
	userRepo := repositories.NewUserRepository(app.dbPool, app.logger)
	shopRepo := repositories.NewShopRepository(app.dbPool, app.logger)
	transactionRepo := repositories.NewTransactionRepository(app.dbPool, app.logger)

	// Инициализация сервисного слоя
	txExecutor := services.NewTxExecutor(app.dbPool, app.logger)
	userService := services.NewUserService(userRepo, shopRepo, transactionRepo, txExecutor, app.logger)
	transactionService := services.NewTransactionService(userRepo, transactionRepo, txExecutor, app.logger)
	shopService := services.NewShopService(userRepo, shopRepo, txExecutor, app.logger)

	// Инициализация обработчиков
	userHandler := delivery.NewUserHandler(userService, token)
	transactionHandler := delivery.NewTransactionHandler(transactionService)
	shopHandler := delivery.NewShopHandler(shopService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(token, secretKey, app.logger)

	// Настройка маршрутов API
	router := gin.Default()
	app.RegisterRoutes(router, userHandler, transactionHandler, shopHandler, authMiddleware)

	// Формируем адрес для сервера из конфигурации
	host := app.config.ApiServerConfig.Host
	port := app.config.ApiServerConfig.Port
	addr := net.JoinHostPort(host, port)
	app.apiServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return nil
}

// newDBConn устанавливает подключение к базе данных с использованием строки подключения
func newDBConn(dbcfg *config.Database) (*pgxpool.Pool, error) {
	// Получаем строку подключения
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		dbcfg.Driver,
		dbcfg.User,
		dbcfg.Password,
		dbcfg.Host,
		dbcfg.DBPort,
		dbcfg.Name,
		dbcfg.SSLMode,
	)

	// Устанавливаем соединение с базой данных
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return nil, err
	}

	return pool, nil
}
