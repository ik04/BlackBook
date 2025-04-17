package model

import "github.com/google/uuid"

type Tag struct{
	ID uuid.UUID `db:"id"`
	UserID string `db:"user_id"`
	Name string `db:"name"`
}