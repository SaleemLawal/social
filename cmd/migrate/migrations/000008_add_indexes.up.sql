CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_comments_content ON COMMENTS USING GIN (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_title ON POSTS USING GIN (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_posts_tags ON POSTS USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_users_username ON USERS (username);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON POSTS (user_id);

CREATE INDEX IF NOT EXISTS idx_comments_post_id ON COMMENTS (post_id);
