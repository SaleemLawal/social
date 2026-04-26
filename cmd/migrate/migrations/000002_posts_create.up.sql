CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    user_id bigint NOT NULL REFERENCES users(id),
    tags text[],
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);