package usecase

import (
	"context"
	"fmt"

	"github.com/igkostyuk/tasktracker/domain"
)

type columnUsecase struct {
	columnRepo domain.ColumnRepository
	taskRepo   domain.TaskRepository
}

// New will create new a ColumnUsecase object representation of domain.ColumnUsecase interface.
func New(c domain.ColumnRepository, t domain.TaskRepository) domain.ColumnUsecase {
	return &columnUsecase{columnRepo: c, taskRepo: t}
}

func (c *columnUsecase) Fetch(ctx context.Context) ([]domain.Column, error) {
	return c.columnRepo.Fetch(ctx)
}

func (c *columnUsecase) FetchTasks(ctx context.Context, id string) ([]domain.Task, error) {
	if _, err := c.columnRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fetch tasks by column id: %w", err)
	}

	return c.taskRepo.FetchByColumnID(ctx, id)
}

func (c *columnUsecase) FetchByProjectID(ctx context.Context, id string) ([]domain.Column, error) {
	return c.columnRepo.FetchByProjectID(ctx, id)
}

func (c *columnUsecase) GetByID(ctx context.Context, id string) (domain.Column, error) {
	return c.columnRepo.GetByID(ctx, id)
}

func (c *columnUsecase) Update(ctx context.Context, pr *domain.Column) error {
	columns, err := c.columnRepo.FetchByProjectID(ctx, pr.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch by project id: %w", err)
	}
	for _, column := range columns {
		if column.Name == pr.Name {
			return domain.ErrColumnName
		}
	}

	return c.columnRepo.Update(ctx, pr)
}

func (c *columnUsecase) Store(ctx context.Context, cm *domain.Column) error {
	columns, err := c.columnRepo.FetchByProjectID(ctx, cm.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch by project id: %w", err)
	}
	for _, column := range columns {
		if column.Name == cm.Name {
			return domain.ErrColumnName
		}
	}

	return c.columnRepo.Store(ctx, cm)
}

func (c *columnUsecase) Delete(ctx context.Context, id string) error {
	column, err := c.columnRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get column by id: %w", err)
	}
	columns, err := c.columnRepo.FetchByProjectID(ctx, column.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch columns by project id: %w", err)
	}
	if len(columns) == 1 {
		return domain.ErrLastColumn
	}
	tasks, err := c.taskRepo.FetchByColumnID(ctx, column.ID)
	if err != nil {
		return fmt.Errorf("fetch tasks by column id: %w", err)
	}
	if len(tasks) == 0 {
		return c.columnRepo.Delete(ctx, id)
	}
	// TODO move tasks and change position
	return c.columnRepo.Delete(ctx, id)
}
