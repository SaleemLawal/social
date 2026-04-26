CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email citext NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

