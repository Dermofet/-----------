ALTER TABLE IF EXISTS release_date DROP COLUMN IF EXISTS release_date;

ALTER TABLE user_music
    DROP FOREIGN KEY user_music_user_id_fkey,
    ADD CONSTRAINT user_music_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users (id);