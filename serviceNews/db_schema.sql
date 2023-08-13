DROP TABLE IF EXISTS posts, sources;

DROP TYPE IF EXISTS sourceFetchUpdates;
CREATE TYPE sourceFetchUpdates AS ENUM ('on', 'off');

CREATE TABLE sources (
    id SERIAL PRIMARY KEY,
    link TEXT NOT NULL,
    fetchUpdates sourceFetchUpdates NOT NULL DEFAULT 'on',
    defaultInterval INTEGER NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    source INTEGER REFERENCES sources(id) NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pubTime BIGINT NOT NULL,
    link TEXT NOT NULL,
    guid TEXT NOT NULL,
    UNIQUE (source, guid)
);
