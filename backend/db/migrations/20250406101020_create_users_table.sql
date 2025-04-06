-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) UNIQUE NOT NULL,
    email_verified_at TIMESTAMP
);

CREATE INDEX idx_users_name ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_name;
DROP INDEX idx_users_email;

DROP TABLE users;
-- +goose StatementEnd
