package repository

import "database/sql"

type RepostRepository struct {
	db *sql.DB
}

func NewRepostRepository(db *sql.DB) *RepostRepository {
	return &RepostRepository{db: db}
}

func (r *RepostRepository) Toggle(userID, postID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM reposts WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		_, err = r.db.Exec("DELETE FROM reposts WHERE user_id = ? AND post_id = ?", userID, postID)
		return false, err
	}
	_, err = r.db.Exec("INSERT INTO reposts (user_id, post_id) VALUES (?, ?)", userID, postID)
	return true, err
}

func (r *RepostRepository) Count(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM reposts WHERE post_id = ?", postID).Scan(&count)
	return count, err
}

func (r *RepostRepository) HasReposted(userID, postID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM reposts WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&count)
	return count > 0, err
}
