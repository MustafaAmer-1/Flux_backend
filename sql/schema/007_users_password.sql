-- +goose Up
ALTER TABLE users ADD COLUMN passwd BYTEA NOT NULL DEFAULT(
    random()::text::bytea
);
ALTER TABLE users ADD CONSTRAINT unique_name UNIQUE (name);


-- +goose Down
ALTER TABLE users DROP COLUMN passwd;
ALTER TABLE users DROP CONSTRAINT unique_name;
