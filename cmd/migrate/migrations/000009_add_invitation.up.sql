CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);