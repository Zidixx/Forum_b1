package repository

import (
	"database/sql"
	"forum/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	result, err := r.db.Exec(
		"INSERT INTO comments (post_id, user_id, parent_id, content) VALUES (?, ?, ?, ?)",
		comment.PostID, comment.UserID, comment.ParentID, comment.Content,
	)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	comment.ID = int(id)
	return nil
}

func (r *CommentRepository) FindByPostID(postID int) ([]models.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.created_at, c.updated_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.ParentID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.Author); err != nil {
			return nil, err
		}
		all = append(all, c)
	}

	// Build nested structure: top-level comments with their replies
	var topLevel []models.Comment
	replyMap := make(map[int][]models.Comment)

	for i := range all {
		if all[i].ParentID == 0 {
			topLevel = append(topLevel, all[i])
		} else {
			replyMap[all[i].ParentID] = append(replyMap[all[i].ParentID], all[i])
		}
	}

	for i := range topLevel {
		if replies, ok := replyMap[topLevel[i].ID]; ok {
			topLevel[i].Replies = replies
		}
	}

	return topLevel, nil
}

func (r *CommentRepository) FindByID(id int) (*models.Comment, error) {
	comment := &models.Comment{}
	err := r.db.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.created_at, c.updated_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?
	`, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.ParentID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.Author)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *CommentRepository) Update(comment *models.Comment) error {
	_, err := r.db.Exec(
		"UPDATE comments SET content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		comment.Content, comment.ID,
	)
	return err
}

func (r *CommentRepository) Delete(id int) error {
	// Delete replies first, then the comment
	r.db.Exec("DELETE FROM comments WHERE parent_id = ?", id)
	_, err := r.db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

func (r *CommentRepository) CountAll() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM comments").Scan(&count)
	return count, err
}
