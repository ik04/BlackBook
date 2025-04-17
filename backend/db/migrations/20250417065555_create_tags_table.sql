-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags (
    id UUID PRIMARY KEY, 
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_tags_user_id ON tags(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_tags_user_id;

DROP TABLE tags;
-- +goose StatementEnd
