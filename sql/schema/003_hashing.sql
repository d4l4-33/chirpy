-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password TEXT;

ALTER TABLE users
ALTER hashed_password SET DEFAULT 'unset';

ALTER TABLE users
ALTER hashed_password SET NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;