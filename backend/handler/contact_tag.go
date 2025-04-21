package handler

import (
	"backend/repository"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ContactTagHandler struct {
	repo        *repository.ContactTagRepository
	contactRepo *repository.ContactRepository
	tagRepo     *repository.TagRepository
}

func NewContactTagHandler(
	ctRepo *repository.ContactTagRepository,
	contactRepo *repository.ContactRepository,
	tagRepo *repository.TagRepository,
) *ContactTagHandler {
	return &ContactTagHandler{
		repo:        ctRepo,
		contactRepo: contactRepo,
		tagRepo:     tagRepo,
	}
}

func (h *ContactTagHandler) AddTagToContact(c echo.Context) error {
	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag ID")
	}

	userID := c.Get("user_id").(string)

	contact, err := h.contactRepo.GetByID(contactID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch contact")
	}
	if contact.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	tag, err := h.tagRepo.GetByID(tagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch tag")
	}
	if tag.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	if err := h.repo.AddTagToContact(contactID, tagID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not add tag to contact")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag added to contact"})
}

func (h *ContactTagHandler) RemoveTagFromContact(c echo.Context) error {
	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag ID")
	}

	userID := c.Get("user_id").(string)

	contact, err := h.contactRepo.GetByID(contactID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch contact")
	}
	if contact.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	if err := h.repo.RemoveTagFromContact(contactID, tagID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not remove tag from contact")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "tag removed from contact"})
}

func (h *ContactTagHandler) GetTagsForContact(c echo.Context) error {
	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	userID := c.Get("user_id").(string)

	contact, err := h.contactRepo.GetByID(contactID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve contact")
	}
	if contact.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	tags, err := h.repo.GetTagsForContact(contactID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve tags")
	}

	return c.JSON(http.StatusOK, tags)
}

func (h *ContactTagHandler) GetContactsForTag(c echo.Context) error {
	tagID, err := uuid.Parse(c.Param("tag_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag ID")
	}

	userID := c.Get("user_id").(string)

	tag, err := h.tagRepo.GetByID(tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "tag not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve tag")
	}
	if tag.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	contacts, err := h.repo.GetContactsForTag(tagID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve contacts")
	}

	return c.JSON(http.StatusOK, contacts)
}

