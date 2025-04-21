package repository

import (
	"backend/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)
type ContactTagRepository struct {
	db *sqlx.DB
}

func NewContactTagRepository(db *sqlx.DB) *ContactTagRepository {
	return &ContactTagRepository{db:db}
}

func (r *ContactTagRepository) AddTagToContact(contactID, tagID uuid.UUID) error {
	_, err := r.db.Exec(`
		INSERT INTO contact_tags (id, contact_id, tag_id)
		VALUES (?, ?, ?)
	`, uuid.New(), contactID, tagID)
	return err
}

func (r *ContactTagRepository) RemoveTagFromContact(contactID, tagID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM contact_tags
		WHERE contact_id = ? AND tag_id = ?
	`, contactID, tagID)
	return err
}

func (r *ContactTagRepository) GetTagsForContact(contactID uuid.UUID) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Select(&tags, `
		SELECT t.*
		FROM tags t
		JOIN contact_tags ct ON t.id = ct.tag_id
		WHERE ct.contact_id = ?
	`, contactID)
	return tags, err
}

func (r *ContactTagRepository) GetContactsForTag(tagID uuid.UUID) ([]model.Contact, error) {
	var contacts []model.Contact
	err := r.db.Select(&contacts, `
		SELECT c.*
		FROM contacts c
		JOIN contact_tags ct ON c.id = ct.contact_id
		WHERE ct.tag_id = ?
	`, tagID)
	return contacts, err
}

func (r *ContactTagRepository) RemoveAllTagsFromContact(contactID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM contact_tags
		WHERE contact_id = ?
	`, contactID)
	return err
}

func (r *ContactTagRepository) RemoveAllContactsFromTag(tagID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM contact_tags
		WHERE tag_id = ?
	`, tagID)
	return err
}
