CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL,
    name TEXT UNIQUE,
    email TEXT UNIQUE,
    password TEXT
);