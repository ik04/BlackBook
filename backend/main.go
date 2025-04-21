package main

import (
	"backend/db"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
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

	contactRepo := repository.NewContactRepository(conn)
	tagRepo := repository.NewTagRepository(conn)
	contactTagRepo := repository.NewContactTagRepository(conn)

	tagHandler := handler.NewTagHandler(tagRepo)
	contactTagHandler := handler.NewContactTagHandler(contactTagRepo, contactRepo, tagRepo)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	a := e.Group("/auth")
	a.POST("/register", userHandler.RegisterUser)
	a.POST("/login", userHandler.LoginUser)
	a.GET("/me", userHandler.FetchUserData)

	api := e.Group("/api")
	api.Use(middleware.JWTMiddleware)

	contacts := api.Group("/contacts")
	contacts.POST("", contactHandler.CreateContact)
	contacts.GET("", contactHandler.GetAllContacts)
	contacts.GET("/:id", contactHandler.GetContactByID)
	contacts.PUT("/:id", contactHandler.UpdateContact)
	contacts.DELETE("/:id", contactHandler.DeleteContact)

	contacts.POST("/:contact_id/tags/:tag_id", contactTagHandler.AddTagToContact)
	contacts.DELETE("/:contact_id/tags/:tag_id", contactTagHandler.RemoveTagFromContact)
	contacts.GET("/:contact_id/tags", contactTagHandler.GetTagsForContact)

	api.GET("/tags/:tag_id/contacts", contactTagHandler.GetContactsForTag)

	tags := api.Group("/tags")
	tags.POST("", tagHandler.CreateTag)
	tags.GET("", tagHandler.GetAllTags)
	tags.PUT("/:id", tagHandler.UpdateTag)
	tags.DELETE("/:id", tagHandler.DeleteTag)

	e.Logger.Fatal(e.Start(os.Getenv("APP_PORT")))
}
