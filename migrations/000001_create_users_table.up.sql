-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    username TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000
);