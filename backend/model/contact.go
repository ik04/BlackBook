package model

import "github.com/google/uuid"

type Contact struct{
	ID uuid.UUID `db:"id"`
	Name string `db:"name"`
	PhoneNumber string `db:"phone_number"`
	ProfilePicture string `db:"profile_picture"`
	Description	string `db:"description"`
	UserID string `db:"user_id"`
}