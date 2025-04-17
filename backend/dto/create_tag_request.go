package dto

import (
	"backend/model"

	"github.com/google/uuid"
)

type CreateTagRequest struct {
	Name           string `json:"name" validate:"required"`
}

func (c *CreateTagRequest) ToModel(userID string) *model.Tag{
	return &model.Tag{
		ID: uuid.New(),
		UserID: userID,
		Name: c.Name,
	}
}