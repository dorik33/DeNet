package taskrepo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/dorik33/DeNet/internal/models"
	"github.com/dorik33/DeNet/internal/repository"
	storeerrors "github.com/dorik33/DeNet/internal/repository/storeErorrs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepository struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewTaskRepository(pool *pgxpool.Pool, log *slog.Logger) repository.TaskRepository {
	return &taskRepository{
		pool: pool,
		log:  log,
	}
}

func (repo *taskRepository) CompleteTask(ctx context.Context, userID int, taskID int) error {
	query := `
		INSERT INTO user_tasks (user_id, task_id, completed_at)
		VALUES ($1, $2, now());
	`

	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("user_id", userID), slog.Int("task_id", taskID))

	_, err := repo.pool.Exec(ctx, query, userID, taskID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return storeerrors.ErrTaskCompleted
			}
		}
		repo.log.Error("Failed to complete task", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (repo *taskRepository) GetTaskByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
        SELECT id, name, description, reward, created_at
        FROM tasks
        WHERE id = $1;
    `

	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("id", id))

	var task models.Task
	err := repo.pool.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Reward,
		&task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storeerrors.ErrTaskNotFound
		}
		repo.log.Error("Failed to get task", slog.String("error", err.Error()))
		return nil, err
	}
	return &task, nil
}

func (repo *taskRepository) GetUserTasks(ctx context.Context, userID int) ([]models.Task, error) {
	query := `
        SELECT t.id, t.name, t.description, t.reward
        FROM tasks t
        INNER JOIN user_tasks ut ON t.id = ut.task_id
        WHERE ut.user_id = $1;
    `

	repo.log.Debug("Executing query", slog.String("query", query), slog.Int("user_id", userID))

	rows, err := repo.pool.Query(ctx, query, userID)
	if err != nil {
		repo.log.Error("Failed to get completed tasks", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Reward); err != nil {
			repo.log.Error("Failed to scan task", slog.String("error", err.Error()))
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
