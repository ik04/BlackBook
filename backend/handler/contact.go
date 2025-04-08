package handler

import "github.com/jmoiron/sqlx"

type ContactHandler struct{
	db *sqlx.DB
}