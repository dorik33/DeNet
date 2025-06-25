package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	HashPassword []byte    `json:"-"`
	ReferrerID   *int      `json:"referrer_id,omitempty"`
	Points       int       `json:"points"`
	CreatedAt    time.Time `json:"created_at"`
}

type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Reward      int       `json:"reward"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserTask struct {
	UserID      int       `json:"user_id"`
	TaskID      int       `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
}

type UserStatus struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Points     int    `json:"points"`
	ReferrerID *int   `json:"referrer_id,omitempty"`
	Tasks      []Task `json:"tasks"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string
}
