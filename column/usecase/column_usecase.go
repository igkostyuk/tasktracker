package usecase

import (
	"context"

	"github.com/igkostyuk/tasktracker/domain"
)

type columnUsecase struct {
	columnRepo domain.ColumnRepository
	taskRepo   domain.TaskRepository
}

func New(c domain.ColumnRepository, t domain.TaskRepository) domain.ColumnUsecase {
	return &columnUsecase{columnRepo: c, taskRepo: t}
}

func (c *columnUsecase) Fetch(ctx context.Context) ([]domain.Column, error) {
	return c.columnRepo.Fetch(ctx)
}

func (c *columnUsecase) FetchTasks(ctx context.Context, id string) ([]domain.Task, error) {
	return c.taskRepo.FetchByColumnID(ctx, id)
}

func (c *columnUsecase) FetchByProjectID(ctx context.Context, id string) ([]domain.Column, error) {
	return c.columnRepo.FetchByProjectID(ctx, id)
}

func (c *columnUsecase) GetByID(ctx context.Context, id string) (domain.Column, error) {
	return c.columnRepo.GetByID(ctx, id)
}

func (c *columnUsecase) Update(ctx context.Context, pr *domain.Column) error {
	return c.columnRepo.Update(ctx, pr)
}

func (c *columnUsecase) Store(ctx context.Context, m *domain.Column) error {
	columns, err := c.FetchByProjectID(ctx, m.ProjectID)
	if err != nil {
		return err
	}
	for _, column := range columns {
		if column.Name == m.Name {
			return domain.ErrColumnName
		}
	}

	return c.columnRepo.Store(ctx, m)
}

func (c *columnUsecase) Delete(ctx context.Context, id string) error {
	column, err := c.columnRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	columns, err := c.FetchByProjectID(ctx, column.ProjectID)
	if err != nil {
		return err
	}
	if len(columns) == 1 {
		return domain.ErrLastColumn
	}
	return c.columnRepo.Delete(ctx, id)
}
