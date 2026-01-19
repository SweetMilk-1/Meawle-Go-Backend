-- Удаление таблиц в обратном порядке (сначала зависимые таблицы)
DROP TABLE IF EXISTS cat_house_occupancy;
DROP TABLE IF EXISTS cat_houses;

-- Удаление индексов
DROP INDEX IF EXISTS idx_cat_houses_name;
DROP INDEX IF EXISTS idx_cat_house_occupancy_cat_id;
DROP INDEX IF EXISTS idx_cat_house_occupancy_cat_house_id;