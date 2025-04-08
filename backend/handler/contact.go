package handler

import (
	"backend/dto"
	"backend/middleware"
	"backend/repository"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type ContactHandler struct {
	repo *repository.ContactRepository
}

func NewContactHandler(db *sqlx.DB) *ContactHandler {
	return &ContactHandler{
		repo: repository.NewContactRepository(db),
	}
}

func (h *ContactHandler) CreateContact(c echo.Context) error {
	var req dto.CreateContactRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	contact := req.ToModel(userID)

	if err := h.repo.Create(contact); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not create contact")
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "contact created"})
}

func (h *ContactHandler) GetAllContacts(c echo.Context) error {
	userID := c.Get("user_id").(string)

	contacts, err := h.repo.GetAllByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve contacts")
	}

	return c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) GetContactByID(c echo.Context) error {
	idParam := c.Param("id")
	contactID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	contact, err := h.repo.GetByID(contactID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve contact")
	}

	userID := c.Get("user_id").(string)
	if contact.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	return c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) UpdateContact(c echo.Context) error {
	idParam := c.Param("id")
	contactID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	var req dto.CreateContactRequest
	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := c.Get("user_id").(string)
	contact := req.ToModel(userID)
	contact.ID = contactID

	existing, err := h.repo.GetByID(contactID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve contact")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	if err := h.repo.Update(contact); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update contact")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "contact updated"})
}

func (h *ContactHandler) DeleteContact(c echo.Context) error {
	idParam := c.Param("id")
	contactID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contact ID")
	}

	userID := c.Get("user_id").(string)

	existing, err := h.repo.GetByID(contactID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "contact not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete contact")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "access denied")
	}

	if err := h.repo.Delete(contactID, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete contact")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "contact deleted"})
}
