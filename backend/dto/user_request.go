package dto

import (
	"backend/model"
	"database/sql"

	"github.com/google/uuid"
)


type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (req *CreateUserRequest) ToModel(hashedPassword string) *model.User {
	return &model.User{
		ID:              uuid.New(),
		Email:           req.Email,
		Username:        req.Username,
		Password:        hashedPassword,
		EmailVerifiedAt: sql.NullTime{Valid: false},
	}
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}