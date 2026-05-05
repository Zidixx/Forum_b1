package models

import "time"

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	ParentID  int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    string
	Likes     int
	Dislikes  int
	UserVote  string
	Replies   []Comment
}
