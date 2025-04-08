-- +goose Up
-- +goose StatementBegin
CREATE TABLE contacts (
    id UUID PRIMARY KEY, 
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    profile_picture VARCHAR(255),
    description VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE


);
CREATE INDEX idx_contacts_phone_number ON contacts(phone_number);
CREATE INDEX idx_contacts_name ON contacts(name);
CREATE INDEX idx_contacts_user_id ON contacts(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_contacts_phone_number;
DROP INDEX idx_contacts_name;

DROP TABLE contacts;
-- +goose StatementEnd
