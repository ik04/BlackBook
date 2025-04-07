package dto

import (
	"backend/model"
)

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewUserResponse(u *model.User) *UserResponse {
	return &UserResponse{
		ID:       u.ID.String(),
		Email:    u.Email,
		Username: u.Username,
	}
}
