-- Создание таблицы с товарами
CREATE TABLE IF NOT EXISTS products (
    item VARCHAR(20) PRIMARY KEY,
    price INT NOT NULL
);

-- Заполнение таблицы товарами
INSERT INTO products (item, price) VALUES
                                   ('t-shirt', 80),
                                   ('cup', 20),
                                   ('book', 50),
                                   ('pen', 10),
                                   ('powerbank', 200),
                                   ('hoody', 300),
                                   ('umbrella', 200),
                                   ('socks', 10),
                                   ('wallet', 50),
                                   ('pink-hoody', 500);

