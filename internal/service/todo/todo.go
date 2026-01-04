package todo

import (
	"TrueToDoList/internal/domain"
	"TrueToDoList/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type ToDo struct {
	log     *slog.Logger
	Storage Storage
}

type Storage interface {
	SaveTask(ctx context.Context, task domain.Task) (int64, error)
	GetTask(ctx context.Context, id *int64) (map[int64]domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	DoneTask(ctx context.Context, id int64) error
}

func New(logger *slog.Logger, storage Storage) *ToDo {
	return &ToDo{
		log:     logger,
		Storage: storage,
	}
}
func (t *ToDo) GetTasks(ctx context.Context, id *int64) (map[int64]domain.Task, error) {
	const op = "service.todo.GetTasks"

	log := t.log.With(
		slog.String("op", op),
	)

	log.Info("getting tasks")

	tasks, err := t.Storage.GetTask(ctx, id)
	if err != nil && id != nil { // TODO fix id
		log.Error("failed to get task")
		if errors.Is(err, storage.ErrNoSuchTask) {
			return nil, storage.ErrNoSuchTask
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("task was received successfully")
	return tasks, nil
}

func (t *ToDo) CreateTask(ctx context.Context, title string, description string) (int64, error) {
	const op = "service.todo.CreateTask"
	log := t.log.With(
		slog.String("op", op),
	)
	log.Info("creating task")
	id, err := t.Storage.SaveTask(ctx, domain.Task{Title: title, Description: description})
	if err != nil {
		log.Error("failed to create task", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("task was created")
	return id, nil
}
func (t *ToDo) DeleteTask(ctx context.Context, id int64) error {
	const op = "service.todo.DeleteTask"

	log := t.log.With(slog.String("op", op))

	log.Info("deleting task")

	if err := t.Storage.DeleteTask(ctx, id); err != nil {
		log.Error("failed to delete task")
		if errors.Is(err, storage.ErrNoSuchTask) {
			return storage.ErrNoSuchTask
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (t *ToDo) DoneTask(ctx context.Context, id int64) error {
	const op = "service.todo.DoneTask"

	log := t.log.With("op", op)

	log.Info("starting to done task")

	if err := t.Storage.DoneTask(ctx, id); err != nil {
		log.Error("failed to done task")
		if errors.Is(err, storage.ErrNoSuchTask) {
			return storage.ErrNoSuchTask
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
