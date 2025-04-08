package handler

import (
	"backend/dto"
	"backend/middleware"
	"backend/repository"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{
	db *sqlx.DB
	userRepo *repository.UserRepository
}

func NewUserHandler(db *sqlx.DB) *UserHandler{
	return &UserHandler{
		db: db,
		userRepo: repository.NewUserRepository(db),
	}
}

func (h *UserHandler) RegisterUser(c echo.Context) error {
	var userRepo = h.userRepo
	var req dto.CreateUserRequest

	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := req.ToModel(string(hashed))
	userRepo.Create(user)
	return c.JSON(http.StatusCreated, echo.Map{"message": "user created"})
}
func (h *UserHandler) LoginUser(c echo.Context) error {
	var userRepo = h.userRepo
	var req dto.LoginRequest

	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := userRepo.GetByLogin(req.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or username")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"username": user.Username,
		"exp":     jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), 
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := []byte(os.Getenv("JWT_SECRET")) 
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not generate token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
		"token":   signedToken,
	})
}

func (h *UserHandler) FetchUserData(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not in context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID format in token")
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "could not fetch user")
	}

	user.Password = ""

	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

