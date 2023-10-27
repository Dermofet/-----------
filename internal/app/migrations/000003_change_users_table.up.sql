ALTER TABLE users
    ADD COLUMN username varchar(255) NOT NULL,
    DROP COLUMN first_name,
    DROP COLUMN second_name,
    DROP COLUMN last_name,
    DROP COLUMN age,
    DROP COLUMN email,
    DROP COLUMN phone;