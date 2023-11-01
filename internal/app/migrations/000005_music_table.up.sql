ALTER TABLE music
    ADD COLUMN release_date date;

ALTER TABLE user_music
    DROP CONSTRAINT user_music_user_id_fkey,
    ADD CONSTRAINT user_music_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;