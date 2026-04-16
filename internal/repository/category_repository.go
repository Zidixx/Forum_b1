package repository

import (
	"database/sql"
	"forum/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, created_at FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) FindByID(id int) (*models.Category, error) {
	cat := &models.Category{}
	err := r.db.QueryRow("SELECT id, name, created_at FROM categories WHERE id = ?", id).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *CategoryRepository) FindByPostID(postID int) ([]models.Category, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.name, c.created_at
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
		ORDER BY c.name
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
