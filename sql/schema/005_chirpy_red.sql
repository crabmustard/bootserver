-- +goose Up
ALTER TABLE users
ADD chirpy_red BOOLEAN DEFAULT FALSE NOT NULL;

-- +goose Down
ALTER TABLE users
DROP COLUMN chirpy_red;