-- +migrate Up
ALTER TABLE users ADD COLUMN is_email_confirmed BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN email_confirmation_token VARCHAR(128);
