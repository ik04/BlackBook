package model

import (
	"database/sql"

	"github.com/google/uuid"
)


type User struct{
	ID uuid.UUID `db:"id"`
	Email string `db:"email"`
	Username string `db:"username"`
	Password string `db:"password"`
	EmailVerifiedAt sql.NullTime `db:"email_verified_at"`
}