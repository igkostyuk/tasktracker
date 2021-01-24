package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out ./mock/task.go -pkg mocks . TaskUsecase TaskRepository

// Task represent a task in tasktracker.
type Task struct {
	ID          uuid.UUID `json:"id" readonly:"true"`
	Position    int       `json:"position" validate:"required,min=0"`
	Name        string    `json:"name" validate:"required,min=1,max=500"`
	Description string    `json:"description" validate:"required,min=0,max=5000"`
	ColumnID    uuid.UUID `json:"column_id" validate:"required"`
}

// TaskUsecase represent the task's usecases.
type TaskUsecase interface {
	Fetch(ctx context.Context) ([]Task, error)
	FetchByColumnID(ctx context.Context, id uuid.UUID) ([]Task, error)
	FetchByProjectID(ctx context.Context, id uuid.UUID) ([]Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (Task, error)
	Update(ctx context.Context, tk *Task) error
	Store(context.Context, *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	FetchComments(ctx context.Context, id uuid.UUID) ([]Comment, error)
}

// TaskRepository represent the project's repository contract.
type TaskRepository interface {
	Fetch(ctx context.Context) ([]Task, error)
	FetchByColumnID(ctx context.Context, id uuid.UUID) ([]Task, error)
	FetchByProjectID(ctx context.Context, id uuid.UUID) ([]Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (Task, error)
	Update(ctx context.Context, tks ...Task) error
	Store(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
