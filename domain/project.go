package domain

import (
	"context"

	"github.com/gofrs/uuid"
)

// Project ...
type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// ProjectUsecase represent the project's usecases.
type ProjectUsecase interface {
	GetByID(ctx context.Context, id int64) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(context.Context, *Project) error
	Delete(ctx context.Context, id int64) error
}

// ProjectRepository represent the project's repository contract.
type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (Project, error)
	Update(ctx context.Context, pr *Project) error
	Store(ctx context.Context, a *Project) error
	Delete(ctx context.Context, id int64) error
}
