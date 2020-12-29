package domain

import (
	"context"
)

// Project ...
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectUsecase represent the project's usecases.
type ProjectUsecase interface {
	Fetch(ctx context.Context) ([]Project, error)
	GetByID(ctx context.Context, id string) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(context.Context, *Project) error
	Delete(ctx context.Context, id string) error
}

// ProjectRepository represent the project's repository contract.
type ProjectRepository interface {
	Fetch(ctx context.Context) ([]Project, error)
	GetByID(ctx context.Context, id string) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(ctx context.Context, a *Project) error
	Delete(ctx context.Context, id string) error
}
