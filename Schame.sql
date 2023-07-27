CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL,
    email TEXT UNIQUE,
    password TEXT
);