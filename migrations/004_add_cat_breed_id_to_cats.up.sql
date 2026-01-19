-- Добавляем поле cat_breed_id в таблицу cats
ALTER TABLE cats ADD COLUMN cat_breed_id INTEGER;

-- В SQLite нельзя добавить FOREIGN KEY через ALTER TABLE,
-- поэтому создаем новую таблицу с внешним ключом и копируем данные

-- Создаем временную таблицу с новой структурой
CREATE TABLE cats_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    description TEXT,
    user_id INTEGER NOT NULL,
    cat_breed_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (cat_breed_id) REFERENCES cat_breeds(id) ON DELETE SET NULL
);

-- Копируем данные из старой таблицы
INSERT INTO cats_new (id, name, age, description, user_id, created_at)
SELECT id, name, age, description, user_id, created_at FROM cats;

-- Удаляем старую таблицу
DROP TABLE cats;

-- Переименовываем новую таблицу
ALTER TABLE cats_new RENAME TO cats;

-- Создаем индексы
CREATE INDEX IF NOT EXISTS idx_cats_user_id ON cats(user_id);
CREATE INDEX IF NOT EXISTS idx_cats_name ON cats(name);
CREATE INDEX IF NOT EXISTS idx_cats_cat_breed_id ON cats(cat_breed_id);

-- Обновляем существующие данные (устанавливаем породу для тестовых котов)
-- Сопоставляем котов с породами по порядку из cat_breeds таблицы
UPDATE cats SET cat_breed_id = 1 WHERE name = 'Мурзик'; -- Сиамская (id: 1)
UPDATE cats SET cat_breed_id = 2 WHERE name = 'Барсик'; -- Мейн-кун (id: 2)
UPDATE cats SET cat_breed_id = 3 WHERE name = 'Васька'; -- Британская короткошерстная (id: 3)
UPDATE cats SET cat_breed_id = 4 WHERE name = 'Рыжик'; -- Сфинкс (id: 4)
UPDATE cats SET cat_breed_id = 5 WHERE name = 'Снежок'; -- Персидская (id: 5)

