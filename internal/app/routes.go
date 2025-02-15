package app

import (
	h "API-Avito-shop/internal/delivery"
	"API-Avito-shop/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (app *App) RegisterRoutes(r *gin.Engine, userHandler *h.UserHandler, coinHandler *h.TransactionHandler, shopHandler *h.ShopHandler, authMiddleware *middleware.AuthMiddleware) {
	users := r.Group("/api")
	{
		users.POST("/auth", userHandler.AuthHandler)
	}

	private := users.Group("/", authMiddleware.AuthMiddleware())

	{
		private.GET("/info", userHandler.InfoHandler)
		private.POST("/sendCoin", coinHandler.SendCoinHandler)
		private.GET("/buy/:item", shopHandler.BuyHandler)
	}
}
