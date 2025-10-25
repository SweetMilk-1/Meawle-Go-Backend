-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    is_admin BIT NOT NULL DEFAULT 0
);

-- Создание индекса для email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Вставка тестовых данных
INSERT INTO users (email, password, is_admin) VALUES 
    ('ivan@example.com', 'admin', 1),
    ('maria@example.com', 'user', 0),
    ('alex@example.com', 'user', 0);