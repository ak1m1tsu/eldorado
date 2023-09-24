package tasks

import (
	"context"
	"fmt"

	"github.com/romankravchuk/eldorado/internal/data"
	"github.com/romankravchuk/eldorado/internal/storages"
	"github.com/romankravchuk/eldorado/internal/storages/tasks"
	"github.com/romankravchuk/eldorado/internal/storages/tasks/pg"
)

type ServiceOption func(*Service) error

func WithTaskStorage(tasks tasks.Storage) ServiceOption {
	return func(s *Service) error {
		s.tasks = tasks
		return nil
	}
}

func WithTaskPostgresStorage(url string) ServiceOption {
	return func(s *Service) error {
		conn, err := storages.NewDBPool("postgres", url)
		if err != nil {
			return err
		}

		tasks, err := pg.New(conn)
		if err != nil {
			return err
		}

		return WithTaskStorage(tasks)(s)
	}
}

type Service struct {
	tasks tasks.Storage
}

func New(opts ...ServiceOption) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Service) List(ctx context.Context, userID string) ([]data.Task, error) {
	const op = "service.task.List"

	tasks, err := s.tasks.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *Service) Create(ctx context.Context, userID string, t data.Task) (data.Task, error) {
	const op = "service.task.Create"

	t.UserID = userID

	if err := s.tasks.Save(ctx, &t); err != nil {
		return data.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	const op = "service.task.Delete"

	if err := s.tasks.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id string, t data.Task) (data.Task, error) {
	const op = "service.task.Update"

	t.ID = id

	if err := s.tasks.Update(ctx, &t); err != nil {
		return data.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}
