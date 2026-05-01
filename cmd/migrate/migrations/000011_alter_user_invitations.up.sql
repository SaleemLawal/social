ALTER TABLE user_invitations 
ADD COLUMN IF NOT EXISTS expires_at timestamp(0) with time zone NOT NULL DEFAULT NOW() + INTERVAL '24 hours';