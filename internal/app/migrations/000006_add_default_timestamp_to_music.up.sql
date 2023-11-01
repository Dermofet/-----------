ALTER TABLE music
ADD COLUMN new_release_date timestamp DEFAULT now();

UPDATE music SET new_release_date = release_date;

ALTER TABLE music
DROP COLUMN release_date;

ALTER TABLE music
RENAME COLUMN new_release_date TO release_date;
