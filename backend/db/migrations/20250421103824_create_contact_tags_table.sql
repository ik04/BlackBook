-- +goose Up
-- +goose StatementBegin
CREATE TABLE contact_tags (
    id UUID PRIMARY KEY, 
    contact_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE contact_tags;
-- +goose StatementEnd
