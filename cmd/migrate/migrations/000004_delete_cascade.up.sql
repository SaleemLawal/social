ALTER TABLE comments
    DROP CONSTRAINT comments_post_id_fkey,
    ADD CONSTRAINT comments_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

ALTER TABLE posts
    DROP CONSTRAINT posts_user_id_fkey,
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

