package dto

import (
	"backend/model"

	"github.com/google/uuid"
)

type CreateContactRequest struct {
	Name           string `json:"name" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	Description    string `json:"description"`
}

func (c *CreateContactRequest) ToModel(userID string) *model.Contact {
	return &model.Contact{
		ID:             uuid.New(),
		Name:           c.Name,
		PhoneNumber:    c.PhoneNumber,
		Description:    c.Description,
		UserID:         userID,
	}
}
