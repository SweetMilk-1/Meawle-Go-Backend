-- Удаляем индекс
DROP INDEX IF EXISTS idx_cats_cat_breed_id;

-- Удаляем внешний ключ
ALTER TABLE cats DROP CONSTRAINT fk_cats_cat_breed_id;

-- Создаем временную таблицу без cat_breed_id
CREATE TABLE cats_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    description TEXT,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Копируем данные из текущей таблицы (без cat_breed_id)
INSERT INTO cats_new (id, name, age, description, user_id, created_at)
SELECT id, name, age, description, user_id, created_at FROM cats;

-- Удаляем текущую таблицу
DROP TABLE cats;

-- Переименовываем новую таблицу
ALTER TABLE cats_new RENAME TO cats;

-- Создаем индексы
CREATE INDEX IF NOT EXISTS idx_cats_user_id ON cats(user_id);
CREATE INDEX IF NOT EXISTS idx_cats_name ON cats(name);
