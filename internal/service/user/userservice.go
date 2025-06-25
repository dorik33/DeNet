package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dorik33/DeNet/internal/config"
	"github.com/dorik33/DeNet/internal/models"
	"github.com/dorik33/DeNet/internal/repository"
	storeerrors "github.com/dorik33/DeNet/internal/repository/storeErorrs"
	"github.com/dorik33/DeNet/internal/service"
	"github.com/dorik33/DeNet/internal/service/serviceerrors"
	"github.com/dorik33/DeNet/internal/utills"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
	taskRepo repository.TaskRepository
	log      *slog.Logger
	cfg      *config.Config
}

func NewUserService(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	log *slog.Logger,
	cfg *config.Config,
) service.UserService {
	return &userService{
		userRepo: userRepo,
		taskRepo: taskRepo,
		log:      log,
		cfg:      cfg,
	}
}

func (service *userService) Register(ctx context.Context, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		service.log.Error("failed to hash password", slog.String("error", err.Error()))
		return fmt.Errorf("failed to hash password: %w", err)
	}
	err = service.userRepo.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		if errors.Is(err, storeerrors.ErrUserExists) {
			return serviceerrors.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	service.log.Info("user created", slog.String("email", email))

	return nil
}

func (service *userService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := service.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storeerrors.ErrUserNotFound) {
			return "", serviceerrors.ErrUserNotFound
		}
		return "", fmt.Errorf("service: failed to get user by email: %w", err)
	}
	err = utills.VerifyPassword(string(user.HashPassword), password)
	if err != nil {
		service.log.Warn("invalid password", slog.String("email", email))
		return "", serviceerrors.ErrInvalidPassword
	}

	token, err := utills.GenerateToken(user.ID, user.Email, []byte(service.cfg.SecretKey), service.cfg.JwtTTL)
	if err != nil {
		service.log.Error("Failed to generate jwt token", slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to generate jwt token: %w", err)
	}

	service.log.Info("User successfully logged", slog.String("email", email))
	return token, nil
}

func (service *userService) GetLeaderboard(ctx context.Context, limit int) ([]models.User, error) {
	users, err := service.userRepo.GetLeaderboard(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	service.log.Info("Leaderboard successfully got")
	return users, nil
}

func (service *userService) SetReferrer(ctx context.Context, userID int, referrerID int) error {
	if referrerID == 0 {
		return serviceerrors.ErrUserNotFound
	}

	err := service.userRepo.SetReferrer(ctx, userID, referrerID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrUserNotFound) {
			return serviceerrors.ErrUserNotFound
		}
		return fmt.Errorf("failed to set referrer: %w", err)
	}

	service.log.Info("Referrer successfully set")
	return nil
}

func (service *userService) Status(ctx context.Context, ID int) (*models.UserStatus, error) {
	if ID == 0 {
		service.log.Error("User not found", slog.Int("userID", ID))
		return nil, serviceerrors.ErrUserNotFound
	}

	user, err := service.userRepo.GetUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrUserNotFound) {
			return nil, serviceerrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userTasks, err := service.taskRepo.GetUserTasks(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tasks: %w", err)
	}
	status := models.UserStatus{
		ID:         user.ID,
		Email:      user.Email,
		Points:     user.Points,
		ReferrerID: user.ReferrerID,
		Tasks:      userTasks,
	}

	service.log.Info("Status successfully got")
	return &status, nil
}

func (service *userService) CompleteTask(ctx context.Context, userID int, taskID int) error {
	task, err := service.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrTaskNotFound) {
			return serviceerrors.ErrTaskNotFound
		}
		return fmt.Errorf("failed to get task: %w", err)
	}

	err = service.taskRepo.CompleteTask(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrTaskCompleted) {
			return serviceerrors.ErrTaskAlreadyDone
		}
		return fmt.Errorf("failed to complete task: %w", err)
	}

	err = service.userRepo.AddPoints(ctx, userID, task.Reward)
	if err != nil {
		return fmt.Errorf("failed to add points: %w", err)
	}

	service.log.Info("Task successfully completed", slog.Int("taskID", taskID), slog.Int("userID", userID), slog.Int("reward", task.Reward))
	return nil
}

