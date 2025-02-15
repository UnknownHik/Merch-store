-- Создание таблицы покупок
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    item VARCHAR(20) NOT NULL,
    price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY (item) REFERENCES products(item) ON DELETE CASCADE
);

-- Добавление индекса для быстрого поиска всех покупок пользователя по его username
CREATE INDEX IF NOT EXISTS idx_username ON purchases(username);