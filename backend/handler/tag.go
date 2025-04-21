package handler

import (
	"backend/middleware"
	"backend/model"
	"backend/repository"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	repo *repository.TagRepository
}

func NewTagHandler(repo *repository.TagRepository) *TagHandler {
	return &TagHandler{repo: repo}
}

func (h *TagHandler) CreateTag(c echo.Context) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	tag := &model.Tag{
		Name:   req.Name,
		UserID: userID,
	}

	if err := h.repo.Create(tag); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create tag")
	}

	return c.JSON(http.StatusCreated, tag)
}

func (h *TagHandler) GetAllTags(c echo.Context) error {
	userID := c.Get("user_id").(string)
	tags, err := h.repo.GetAllByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch tags")
	}
	return c.JSON(http.StatusOK, tags)
}

func (h *TagHandler) UpdateTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag ID")
	}

	var req struct {
		Name string `json:"name" validate:"required"`
	}
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	existing, err := h.repo.GetByID(tagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "tag not found")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	existing.Name = req.Name
	if err := h.repo.Update(existing); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update tag")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag updated"})
}

func (h *TagHandler) DeleteTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag ID")
	}

	userID := c.Get("user_id").(string)
	if err := h.repo.Delete(tagID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete tag")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag deleted"})
}
