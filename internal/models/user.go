package models

import "time"

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
	Role         string
	FavoriteTeam string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
