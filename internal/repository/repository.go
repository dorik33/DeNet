package repository

import (
	"context"

	"github.com/dorik33/DeNet/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email string, password []byte) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SetReferrer(ctx context.Context, userID int, referrerID int) error
	GetLeaderboard(ctx context.Context, limit int) ([]models.User, error)
	AddPoints(ctx context.Context, userID int, points int) error
}

type TaskRepository interface {
	CompleteTask(ctx context.Context, userID int, taskID int) error
	GetTaskByID(ctx context.Context, id int) (*models.Task, error)
	GetUserTasks(ctx context.Context, userID int) ([]models.Task, error)
}
