-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users (
 id SERIAL PRIMARY KEY,
 name VARCHAR(50) NOT NULL,
 surname VARCHAR(50) NOT NULL,
 patronimic VARCHAR(50),
 age INT NOT NULL,
 gender VARCHAR(6) NOT NULL,
 nation VARCHAR(10) NOT NULL

);
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
SELECT 'down SQL query';
-- +goose StatementEnd
