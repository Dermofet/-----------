-- Удаляем столбец "username", если он существует
ALTER TABLE IF EXISTS users DROP COLUMN IF EXISTS username;

-- Добавляем новые столбцы
ALTER TABLE users ADD COLUMN first_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN second_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN age INT NOT NULL;
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
ALTER TABLE users ADD COLUMN phone VARCHAR(255) NOT NULL;
