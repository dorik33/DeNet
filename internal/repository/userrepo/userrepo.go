package userrepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/dorik33/DeNet/internal/models"
	"github.com/dorik33/DeNet/internal/repository"
	storeerrors "github.com/dorik33/DeNet/internal/repository/storeErorrs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, log *slog.Logger) repository.UserRepository {
	return &userRepository{
		pool: pool,
		log:  log,
	}
}

func (repo *userRepository) CreateUser(ctx context.Context, email string, password []byte) error {
	query := `
	INSERT INTO users(email, hash_password) 
	VALUES($1, $2);
	`
	repo.log.Debug("Executing query", slog.String("query", query), slog.String("email", email))
	_, err := repo.pool.Exec(ctx, query, email, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				repo.log.Warn("User with email already exists", slog.String("email", email))
				return storeerrors.ErrUserExists
			}
		}

		repo.log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
	SELECT id, email, hash_password, referrer_id, points, created_at
	FROM users
	WHERE id = $1;
	`
	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("id", id))

	var user models.User
	err := repo.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.HashPassword, &user.ReferrerID, &user.Points, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storeerrors.ErrUserNotFound
		}
		repo.log.Error("Failed to get user", slog.String("error", err.Error()))

		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT id, email, hash_password, referrer_id, points, created_at
	FROM users
	WHERE email = $1;
	`
	repo.log.Debug("Executing query", slog.String("query", query), slog.String("email", email))

	var user models.User
	err := repo.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.HashPassword, &user.ReferrerID, &user.Points, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storeerrors.ErrUserNotFound
		}
		repo.log.Error("Failed to get user", slog.String("error", err.Error()))

		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) SetReferrer(ctx context.Context, userID int, referrerID int) error {
	query := `
		UPDATE users
		SET referrer_id = $1
		WHERE id = $2;
	`

	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("referrer_id", referrerID), slog.Int("user_id", userID))

	cmdTag, err := repo.pool.Exec(ctx, query, referrerID, userID)
	if err != nil {
		repo.log.Error("Failed to set referrer", slog.String("error", err.Error()))
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return storeerrors.ErrUserNotFound
	}

	return nil
}

func (repo *userRepository) GetLeaderboard(ctx context.Context, limit int) ([]models.User, error) {
	query := `
	SELECT id, email, referrer_id, points, created_at
	FROM users
	ORDER BY points DESC
	LIMIT $1;
	`

	repo.log.Debug("Executing query", slog.String("query", query))

	rows, err := repo.pool.Query(ctx, query, limit)
	if err != nil {
		repo.log.Error("Failed to get leaderboard", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.ReferrerID, &user.Points, &user.CreatedAt)
		if err != nil {
			repo.log.Error("Failed to scan user", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) AddPoints(ctx context.Context, userID int, points int) error {
	query := `
        UPDATE users
        SET points = points + $1
        WHERE id = $2;
    `

	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("user_id", userID), slog.Int("points", points))

	cmdTag, err := repo.pool.Exec(ctx, query, points, userID)
	if err != nil {
		repo.log.Error("Failed to add points", slog.String("error", err.Error()))
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return storeerrors.ErrUserNotFound
	}

	return nil
}
