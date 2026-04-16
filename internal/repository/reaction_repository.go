package repository

import (
	"database/sql"
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

// Post reactions

func (r *ReactionRepository) GetPostReaction(userID, postID int) (string, error) {
	var reactionType string
	err := r.db.QueryRow(
		"SELECT type FROM post_reactions WHERE user_id = ? AND post_id = ?",
		userID, postID,
	).Scan(&reactionType)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return reactionType, err
}

func (r *ReactionRepository) SetPostReaction(userID, postID int, reactionType string) error {
	existing, err := r.GetPostReaction(userID, postID)
	if err != nil {
		return err
	}

	if existing == reactionType {
		_, err = r.db.Exec("DELETE FROM post_reactions WHERE user_id = ? AND post_id = ?", userID, postID)
		return err
	}

	if existing != "" {
		_, err = r.db.Exec("UPDATE post_reactions SET type = ? WHERE user_id = ? AND post_id = ?", reactionType, userID, postID)
		return err
	}

	_, err = r.db.Exec("INSERT INTO post_reactions (user_id, post_id, type) VALUES (?, ?, ?)", userID, postID, reactionType)
	return err
}

func (r *ReactionRepository) CountPostReactions(postID int) (likes int, dislikes int, err error) {
	err = r.db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = 'like'", postID).Scan(&likes)
	if err != nil {
		return
	}
	err = r.db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = 'dislike'", postID).Scan(&dislikes)
	return
}

// Comment reactions

func (r *ReactionRepository) GetCommentReaction(userID, commentID int) (string, error) {
	var reactionType string
	err := r.db.QueryRow(
		"SELECT type FROM comment_reactions WHERE user_id = ? AND comment_id = ?",
		userID, commentID,
	).Scan(&reactionType)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return reactionType, err
}

func (r *ReactionRepository) SetCommentReaction(userID, commentID int, reactionType string) error {
	existing, err := r.GetCommentReaction(userID, commentID)
	if err != nil {
		return err
	}

	if existing == reactionType {
		_, err = r.db.Exec("DELETE FROM comment_reactions WHERE user_id = ? AND comment_id = ?", userID, commentID)
		return err
	}

	if existing != "" {
		_, err = r.db.Exec("UPDATE comment_reactions SET type = ? WHERE user_id = ? AND comment_id = ?", reactionType, userID, commentID)
		return err
	}

	_, err = r.db.Exec("INSERT INTO comment_reactions (user_id, comment_id, type) VALUES (?, ?, ?)", userID, commentID, reactionType)
	return err
}

func (r *ReactionRepository) CountCommentReactions(commentID int) (likes int, dislikes int, err error) {
	err = r.db.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND type = 'like'", commentID).Scan(&likes)
	if err != nil {
		return
	}
	err = r.db.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND type = 'dislike'", commentID).Scan(&dislikes)
	return
}
