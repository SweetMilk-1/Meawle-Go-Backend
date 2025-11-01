CREATE TABLE IF NOT EXISTS cats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    description TEXT,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создаем индекс для быстрого поиска по пользователю
CREATE INDEX IF NOT EXISTS idx_cats_user_id ON cats(user_id);

-- Создаем индекс для быстрого поиска по имени
CREATE INDEX IF NOT EXISTS idx_cats_name ON cats(name);

-- Вставляем тестовые данные
INSERT INTO cats (name, age, description, user_id) VALUES
('Мурзик', 3, 'Ласковый и игривый кот', 1),
('Барсик', 5, 'Спокойный и мудрый кот', 1),
('Васька', 2, 'Любит поесть и поспать', 2),
('Рыжик', 4, 'Очень активный и любопытный', 2),
('Снежок', 1, 'Белый пушистый котенок', 3);