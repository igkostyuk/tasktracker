package domain

import (
	"context"
)

//go:generate moq -out ./mock/task.go -pkg mocks . TaskUsecase TaskRepository

// Task represent a task in tasktracker.
type Task struct {
	ID          string `json:"id"`
	Position    int64  `json:"position"`
	Name        string `json:"name" validate:"required,min=1,max=500"`
	Description string `json:"description" validate:"required,min=0,max=5000"`
	ColumnID    string `json:"column_id" validate:"required,uuid4"`
}

// TaskUsecase represent the task's usecases.
type TaskUsecase interface {
	Fetch(ctx context.Context) ([]Task, error)
	FetchByColumnID(ctx context.Context, id string) ([]Task, error)
	FetchByProjectID(ctx context.Context, id string) ([]Task, error)
	GetByID(ctx context.Context, id string) (Task, error)
	Update(ctx context.Context, tk *Task) error
	Store(context.Context, *Task) error
	Delete(ctx context.Context, id string) error
	FetchComments(ctx context.Context, id string) ([]Comment, error)
}

// TaskRepository represent the project's repository contract.
type TaskRepository interface {
	Fetch(ctx context.Context) ([]Task, error)
	FetchByColumnID(ctx context.Context, id string) ([]Task, error)
	FetchByProjectID(ctx context.Context, id string) ([]Task, error)
	GetByID(ctx context.Context, id string) (Task, error)
	Update(ctx context.Context, tk *Task) error
	Store(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id string) error
}
