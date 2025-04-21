package repository

import (
	"backend/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(db *sqlx.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	if tag.ID == uuid.Nil {
		tag.ID = uuid.New()
	}
	_, err := r.db.NamedExec(`
		INSERT INTO tags (id, user_id, name)
		VALUES (:id, :user_id, :name)
	`, tag)
	return err
}

func (r *TagRepository) GetByID(tagID uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Get(&tag, `
		SELECT * FROM tags WHERE id = ?
	`, tagID)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) GetAllByUserID(userID string) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Select(&tags, `
		SELECT id, name FROM tags WHERE user_id = ?
	`, userID)
	return tags, err
}

func (r *TagRepository) Update(tag *model.Tag) error {
	_, err := r.db.Exec(`
		UPDATE tags SET name = ? WHERE id = ? AND user_id = ?
	`, tag.Name, tag.ID, tag.UserID)
	return err
}

func (r *TagRepository) Delete(tagID uuid.UUID, userID string) error {
	_, err := r.db.Exec(`
		DELETE FROM tags WHERE id = ? AND user_id = ?
	`, tagID, userID)
	return err
}
