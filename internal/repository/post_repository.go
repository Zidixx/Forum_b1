package repository

import (
	"database/sql"
	"forum/internal/models"
	"strings"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO posts (user_id, title, content, image_path) VALUES (?, ?, ?, ?)",
		post.UserID, post.Title, post.Content, post.ImagePath,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	post.ID = int(id)

	for _, catID := range categoryIDs {
		if _, err := tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", post.ID, catID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostRepository) FindByID(id int) (*models.Post, error) {
	post := &models.Post{}
	var imagePath sql.NullString
	err := r.db.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &imagePath, &post.CreatedAt, &post.UpdatedAt, &post.Author)
	if err != nil {
		return nil, err
	}
	if imagePath.Valid {
		post.ImagePath = imagePath.String
	}
	return post, nil
}

func (r *PostRepository) FindAll() ([]models.Post, error) {
	return r.queryPosts(`
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
}

func (r *PostRepository) FindByCategory(categoryID int) ([]models.Post, error) {
	return r.queryPosts(`
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC
	`, categoryID)
}

func (r *PostRepository) FindByUserID(userID int) ([]models.Post, error) {
	return r.queryPosts(`
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
		ORDER BY p.created_at DESC
	`, userID)
}

func (r *PostRepository) FindLikedByUserID(userID int) ([]models.Post, error) {
	return r.queryPosts(`
		SELECT p.id, p.user_id, p.title, p.content, p.image_path, p.created_at, p.updated_at, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_reactions pr ON p.id = pr.post_id
		WHERE pr.user_id = ? AND pr.type = 'like'
		ORDER BY p.created_at DESC
	`, userID)
}

func (r *PostRepository) Update(post *models.Post, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE posts SET title = ?, content = ?, image_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		post.Title, post.Content, post.ImagePath, post.ID,
	)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("DELETE FROM post_categories WHERE post_id = ?", post.ID); err != nil {
		return err
	}

	for _, catID := range categoryIDs {
		if _, err := tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", post.ID, catID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE id = ?", id)
	return err
}

func (r *PostRepository) GetCategoryIDs(postID int) ([]int, error) {
	rows, err := r.db.Query("SELECT category_id FROM post_categories WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PostRepository) queryPosts(query string, args ...interface{}) ([]models.Post, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var imagePath sql.NullString
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &imagePath, &p.CreatedAt, &p.UpdatedAt, &p.Author); err != nil {
			return nil, err
		}
		if imagePath.Valid {
			p.ImagePath = imagePath.String
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) Excerpt(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	truncated := content[:maxLen]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}
