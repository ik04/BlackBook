package repository

import (
	"backend/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	_, err := r.db.NamedExec(`
		INSERT INTO users (id, email, username, password, email_verified_at)
		VALUES (:id, :email, :username, :password, :email_verified_at)
	`, user)

	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = ?", id.String())
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	err := r.db.Select(&users, "SELECT * FROM users")
	return users, err
}
