package models

import "time"

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	ImagePath    string
	League       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Author       string
	Categories   []Category
	Likes        int
	Dislikes     int
	UserVote     string // "like", "dislike", or ""
	Comments     []Comment
	Reposts      int
	UserReposted bool
	CommentCount int
}
