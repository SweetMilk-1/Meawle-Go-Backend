-- Создание таблицы пород кошек
CREATE TABLE IF NOT EXISTS cat_breeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_cat_breeds_name ON cat_breeds(name);
CREATE INDEX IF NOT EXISTS idx_cat_breeds_user_id ON cat_breeds(user_id);
CREATE INDEX IF NOT EXISTS idx_cat_breeds_created_at ON cat_breeds(created_at);

-- Вставка тестовых данных
INSERT INTO cat_breeds (name, description, user_id) VALUES 
    ('Сиамская', 'Элегантная кошка с голубыми глазами и характерным окрасом', 1),
    ('Мейн-кун', 'Крупная порода с длинной шерстью и дружелюбным характером', 1),
    ('Британская короткошерстная', 'Крепкая кошка с плюшевой шерстью и спокойным нравом', 2),
    ('Сфинкс', 'Бесшерстная порода с морщинистой кожей и теплым телом', 2),
    ('Персидская', 'Длинношерстная кошка с плоской мордой и спокойным характером', 3);