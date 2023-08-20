DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    id_post BIGINT NOT NULL,
    content TEXT NOT NULL,
    pubTime BIGINT NOT NULL,
    flag_obscene BOOLEAN NOT NULL
);
