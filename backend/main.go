package main

import (
	"backend/db"
	"backend/handler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
    }
	conn := db.Init()
	userHandler := handler.NewUserHandler(conn)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	a := e.Group("/auth")
	a.POST("/register",userHandler.RegisterUser)
	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}