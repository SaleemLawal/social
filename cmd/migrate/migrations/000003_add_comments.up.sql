CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY,
    content TEXT NOT NULL,
    post_id bigserial NOT NULL REFERENCES posts(id),
    user_id bigserial NOT NULL REFERENCES users(id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);