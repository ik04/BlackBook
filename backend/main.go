package main

import (
	"backend/db"
	"backend/handler"
	"backend/middleware"
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
	contactHandler := handler.NewContactHandler(conn)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	a := e.Group("/auth")
	a.POST("/register",userHandler.RegisterUser)
	a.POST("/login",userHandler.LoginUser)

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)
	api.GET("/me",userHandler.FetchUserData)

	contacts := api.Group("/contacts")
	contacts.POST("", contactHandler.CreateContact)
	contacts.GET("", contactHandler.GetAllContacts)
	contacts.GET("/:id", contactHandler.GetContactByID)
	contacts.PUT("/:id", contactHandler.UpdateContact)
	contacts.DELETE("/:id", contactHandler.DeleteContact)


	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}