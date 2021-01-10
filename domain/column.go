package domain

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out ./mock/column.go -pkg mocks . ColumnUsecase ColumnRepository

// Column represent a columns in tasktracker.
type Column struct {
	ID        uuid.UUID `json:"id" readonly:"true"`
	Position  int64     `json:"position"`
	Name      string    `json:"name" validate:"required,min=1,max=255"`
	Status    string    `json:"status" validate:"required,min=1,max=255"`
	ProjectID uuid.UUID `json:"project_id" validate:"required,uuid4"`
}

// ColumnUsecase represent the column's usecases.
type ColumnUsecase interface {
	Fetch(ctx context.Context) ([]Column, error)
	FetchByProjectID(ctx context.Context, id uuid.UUID) ([]Column, error)
	GetByID(ctx context.Context, id uuid.UUID) (Column, error)
	Update(ctx context.Context, cl *Column) error
	Store(context.Context, *Column) error
	Delete(ctx context.Context, id uuid.UUID) error
	FetchTasks(ctx context.Context, id uuid.UUID) ([]Task, error)
}

// ColumnRepository represent the column's repository contract.
type ColumnRepository interface {
	Fetch(ctx context.Context) ([]Column, error)
	FetchByProjectID(ctx context.Context, id uuid.UUID) ([]Column, error)
	GetByID(ctx context.Context, id uuid.UUID) (Column, error)
	Update(ctx context.Context, cl *Column) error
	Store(ctx context.Context, c *Column) error
	Delete(ctx context.Context, id uuid.UUID) error
}
