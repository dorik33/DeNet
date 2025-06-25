package service

import (
	"context"

	"github.com/dorik33/DeNet/internal/models"
)

type UserService interface {
	Register(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (string, error)
	GetLeaderboard(ctx context.Context, limit int) ([]models.User, error)
	SetReferrer(ctx context.Context, userID int, referrerID int) error
	Status(ctx context.Context, ID int) (*models.UserStatus, error)
	CompleteTask(ctx context.Context, userID int, taskID int) error
}
