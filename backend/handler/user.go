package handler

import (
	"backend/dto"
	"backend/middleware"
	"backend/repository"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{
	db *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler{
	return &UserHandler{
		db: db,
	}
}

func (u *UserHandler) RegisterUser(c echo.Context) error {
	var userRepo = repository.NewUserRepository(u.db)
	var req dto.CreateUserRequest

	if err := middleware.BindAndValidate(c, &req); err != nil {
		return err
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := req.ToModel(string(hashed))
	userRepo.Create(user)
	return c.JSON(http.StatusCreated, echo.Map{"message": "user created"})
}