package models

import "time"

type PostReaction struct {
	ID        int
	UserID    int
	PostID    int
	Type      string
	CreatedAt time.Time
}

type CommentReaction struct {
	ID        int
	UserID    int
	CommentID int
	Type      string
	CreatedAt time.Time
}
