package postgresql

import (
	"TrueToDoList/internal/domain"
	"TrueToDoList/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, pool *pgxpool.Pool) (*Storage, error) {
	const op = "storage.postgresql.New"
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		pool: pool,
	}, nil
}

func (s *Storage) SaveTask(ctx context.Context, task domain.Task) (int64, error) {
	const op = "storage.postgresql.SaveTask"
	var id int64
	err := s.pool.QueryRow(ctx, "INSERT INTO tasks (title,description) VALUES ($1,$2) RETURNING id", task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetTask(ctx context.Context, id *int64) (map[int64]domain.Task, error) {
	const op = "storage.postgresql.GetTask"
	if id != nil {
		var idTask int64
		var title string
		var description string
		var done bool
		err := s.pool.QueryRow(ctx, "SELECT id,title,description,done FROM tasks WHERE id = $1", *id).Scan(&idTask, &title, &description, &done)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, storage.ErrNoSuchTask
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		data := make(map[int64]domain.Task)
		data[idTask] = domain.Task{Title: title, Description: description, Done: done}
		return data, nil
	} else {
		rows, err := s.pool.Query(ctx, "SELECT id,title,description,done FROM tasks")
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, "no tasks")
			}
			return nil, err
		}
		defer rows.Close()
		data := make(map[int64]domain.Task)
		for rows.Next() {
			var idTask int64
			var title string
			var description string
			var done bool
			if err := rows.Scan(&idTask, &title, &description, &done); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			data[idTask] = domain.Task{Title: title, Description: description, Done: done}
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return data, nil
	}
}

func (s *Storage) DeleteTask(ctx context.Context, id int64) error {
	const op = "storage.postgresql.DeleteTask"
	cmdTag, err := s.pool.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return storage.ErrNoSuchTask
	}
	return nil
}

func (s *Storage) DoneTask(ctx context.Context, id int64) error {
	const op = "storage.postgresql.DoneTask"
	cmdTag, err := s.pool.Exec(ctx, "UPDATE tasks SET done = TRUE WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return storage.ErrNoSuchTask
	}
	return nil
}
