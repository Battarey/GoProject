-- +migrate Up
ALTER TABLE users ADD COLUMN password_reset_token VARCHAR(128);
ALTER TABLE users ADD COLUMN password_reset_expires_at BIGINT;
