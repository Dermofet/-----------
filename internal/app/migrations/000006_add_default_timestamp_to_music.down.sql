-- Вернуть тип данных столбца release_date к предыдущему
-- Например, если предыдущий тип был date
ALTER TABLE music
ALTER COLUMN release_date TYPE date;

-- Если необходимо вернуть оригинальное значение (например, date)
-- обратно в этот столбец, вы также можете выполнить обратное обновление данных.
-- Например:
-- UPDATE music SET release_date = new_release_date;
