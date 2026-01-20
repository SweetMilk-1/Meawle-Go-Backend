-- Создание таблицы домиков для кошек
CREATE TABLE IF NOT EXISTS cat_houses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0),
    current_occupancy INTEGER NOT NULL DEFAULT 0 CHECK (current_occupancy >= 0 AND current_occupancy <= capacity),
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создание таблицы для связи котов с домиками
CREATE TABLE IF NOT EXISTS cat_house_occupancy (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cat_id INTEGER NOT NULL,
    cat_house_id INTEGER NOT NULL,
    occupied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cat_id) REFERENCES cats(id) ON DELETE CASCADE,
    FOREIGN KEY (cat_house_id) REFERENCES cat_houses(id) ON DELETE CASCADE,
    UNIQUE(cat_id) -- Один кот может быть только в одном домике
);

-- Создание индексов для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_cat_houses_name ON cat_houses(name);
CREATE INDEX IF NOT EXISTS idx_cat_houses_user_id ON cat_houses(user_id);
CREATE INDEX IF NOT EXISTS idx_cat_house_occupancy_cat_id ON cat_house_occupancy(cat_id);
CREATE INDEX IF NOT EXISTS idx_cat_house_occupancy_cat_house_id ON cat_house_occupancy(cat_house_id);

-- Вставляем тестовые данные домиков
INSERT INTO cat_houses (name, capacity, current_occupancy, user_id) VALUES
('Большой дом', 5, 0, 1),
('Средний дом', 3, 0, 1),
('Маленький домик', 1, 0, 2),
('Люкс-апартаменты', 2, 0, 3);