-- +migrate Down
ALTER TABLE users DROP COLUMN is_email_confirmed;
ALTER TABLE users DROP COLUMN email_confirmation_token;
