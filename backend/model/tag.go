package model

import "github.com/google/uuid"

type Tag struct{
	ID uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	Name string `db:"name"`
}