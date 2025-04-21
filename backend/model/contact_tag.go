package model

import "github.com/google/uuid"

type ContactTag struct {
	ID uuid.UUID `db:"id"`
	ContactID string `db:"contact_id"`
	TagID string `db:"tag_id"`
}