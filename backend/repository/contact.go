package repository

import (
	"backend/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ContactRepository struct {
	db *sqlx.DB
}

func NewContactRepository(db *sqlx.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Create(contact *model.Contact) error {
	if contact.ID == uuid.Nil {
		contact.ID = uuid.New()
	}
	_, err := r.db.NamedExec(`
		INSERT INTO contacts (id, name, phone_number, profile_picture, description, user_id)
		VALUES (:id, :name, :phone_number, :profile_picture, :description, :user_id)
	`, contact)
	return err
}

func (r *ContactRepository) GetByID(id uuid.UUID) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.Get(&contact, "SELECT * FROM contacts WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *ContactRepository) GetAllByUserID(userID string) ([]model.Contact, error) {
	var contacts []model.Contact
	err := r.db.Select(&contacts, "SELECT * FROM contacts WHERE user_id = ?", userID)
	return contacts, err
}

func (r *ContactRepository) Delete(id uuid.UUID, userID string) error {
	_, err := r.db.Exec("DELETE FROM contacts WHERE id = ? AND user_id = ?", id, userID)
	return err
}

func (r *ContactRepository) Update(contact *model.Contact) error {
	_, err := r.db.NamedExec(`
		UPDATE contacts
		SET name = :name,
			phone_number = :phone_number,
			profile_picture = :profile_picture,
			description = :description
		WHERE id = :id
	`, contact)
	return err
}
