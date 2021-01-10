package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out ./mock/project.go -pkg mocks . ProjectUsecase ProjectRepository

// Project represent a project in tasktracker.
type Project struct {
	ID          uuid.UUID `json:"id,omitempty" readonly:"true"`
	Name        string    `json:"name" validate:"required,min=1,max=500"`
	Description string    `json:"description" validate:"required,min=0,max=1000"`
}

// ProjectUsecase represent the project's usecases.
type ProjectUsecase interface {
	Fetch(ctx context.Context) ([]Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(context.Context, *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	FetchColumns(ctx context.Context, id uuid.UUID) ([]Column, error)
	FetchTasks(ctx context.Context, id uuid.UUID) ([]Task, error)
}

// ProjectRepository represent the project's repository contract.
type ProjectRepository interface {
	Fetch(ctx context.Context) ([]Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(ctx context.Context, a *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}
