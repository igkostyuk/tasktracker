package usecase

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
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

func (c *columnUsecase) FetchTasks(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
	if _, err := c.columnRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fetch tasks by column id: %w", err)
	}

	return c.taskRepo.FetchByColumnID(ctx, id)
}

func (c *columnUsecase) FetchByProjectID(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
	return c.columnRepo.FetchByProjectID(ctx, id)
}

func (c *columnUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Column, error) {
	return c.columnRepo.GetByID(ctx, id)
}

func (c *columnUsecase) Update(ctx context.Context, cl *domain.Column) error {
	old, err := c.columnRepo.GetByID(ctx, cl.ID)
	if err != nil {
		return fmt.Errorf("fetch by id: %w", err)
	}
	cl.ProjectID = old.ProjectID
	if reflect.DeepEqual(old, *cl) {
		return nil
	}
	columns, err := c.columnRepo.FetchByProjectID(ctx, cl.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch by project id: %w", err)
	}
	if !isNameUnique(columns, cl) {
		return domain.ErrColumnName
	}
	if cl.Position > len(columns)-1 {
		cl.Position = len(columns) - 1
		if reflect.DeepEqual(old, *cl) {
			return nil
		}
	}
	if old.Position < cl.Position {
		return c.MoveRight(ctx, &old, cl, columns)
	}
	if old.Position > cl.Position {
		return c.MoveLeft(ctx, &old, cl, columns)
	}

	return c.columnRepo.Update(ctx, *cl)
}

func (c *columnUsecase) MoveRight(ctx context.Context, old, cl *domain.Column, cls []domain.Column) error {
	cls = cls[old.Position+1 : cl.Position+1]
	for i := range cls {
		cls[i].Position--
	}
	cls = append(cls, *cl)

	return c.columnRepo.Update(ctx, cls...)
}

func (c *columnUsecase) MoveLeft(ctx context.Context, old, cl *domain.Column, cls []domain.Column) error {
	cls = cls[cl.Position:old.Position]
	for i := range cls {
		cls[i].Position++
	}
	cls = append(cls, *cl)

	return c.columnRepo.Update(ctx, cls...)
}

func (c *columnUsecase) Delete(ctx context.Context, id uuid.UUID) error {
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

	leftColumnIndex := column.Position - 1
	if leftColumnIndex < 0 {
		leftColumnIndex = 1
	}
	if err := c.moveTasks(ctx, column.ID, columns[leftColumnIndex].ID); err != nil {
		return err
	}

	return c.columnRepo.Delete(ctx, id)
}

func (c *columnUsecase) moveTasks(ctx context.Context, columnID, leftColumnID uuid.UUID) error {
	tasks, err := c.taskRepo.FetchByColumnID(ctx, columnID)
	if err != nil {
		return fmt.Errorf("fetch tasks by column id: %w", err)
	}
	if len(tasks) == 0 {
		return nil
	}
	leftTasks, err := c.taskRepo.FetchByColumnID(ctx, leftColumnID)
	if err != nil {
		return fmt.Errorf("fetch tasks by left column id: %w", err)
	}
	for i := range tasks {
		tasks[i].ColumnID = leftColumnID
		tasks[i].Position += len(leftTasks)
	}

	return c.taskRepo.Update(ctx, tasks...)
}

func isNameUnique(columns []domain.Column, column *domain.Column) bool {
	for _, c := range columns {
		if column.Name == c.Name && column.ID != c.ID {
			return false
		}
	}

	return true
}
