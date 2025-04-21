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
	cookie, err := c.Cookie("at")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing access token cookie")
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok || userIDStr == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not fetch user")
	}

	user.Password = ""
	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
