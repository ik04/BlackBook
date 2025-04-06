package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)


func Init() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	conn, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to SQLite database")
	return conn
}