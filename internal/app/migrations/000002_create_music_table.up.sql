CREATE TABLE IF NOT EXISTS music (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS user_music (
    user_id UUID NOT NULL,
    music_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (music_id) REFERENCES music (id),
    PRIMARY KEY (user_id, music_id)
);
