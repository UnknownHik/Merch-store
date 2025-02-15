-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_username TEXT NOT NULL,
    to_username TEXT NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_username) REFERENCES users(username) ON DELETE RESTRICT,
    FOREIGN KEY (to_username) REFERENCES users(username) ON DELETE RESTRICT
);

-- Добавление индекса на поле 'from_username' для быстрого поиска от кого получены монеты
CREATE INDEX IF NOT EXISTS idx_from_username ON transactions(from_username);

-- Добавление индекса на поле 'to_username' для быстрого поиска кому отправлены монеты
CREATE INDEX IF NOT EXISTS idx_to_username ON transactions(to_username);