-- +migrate Down
ALTER TABLE users DROP COLUMN password_reset_token;
ALTER TABLE users DROP COLUMN password_reset_expires_at;
