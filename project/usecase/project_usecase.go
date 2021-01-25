package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type projectUsecase struct {
	projectRepo domain.ProjectRepository
	columnRepo  domain.ColumnRepository
	taskRepo    domain.TaskRepository
}

// New will create new a projectUsecase object representation of domain.ProjectUsecase interface.
func New(p domain.ProjectRepository, c domain.ColumnRepository, t domain.TaskRepository) domain.ProjectUsecase {
	return &projectUsecase{projectRepo: p, columnRepo: c, taskRepo: t}
}

func (p *projectUsecase) Fetch(ctx context.Context) ([]domain.Project, error) {
	return p.projectRepo.Fetch(ctx)
}

func (p *projectUsecase) FetchColumns(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
	if _, err := p.projectRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fetch columns by project id: %w", err)
	}

	return p.columnRepo.FetchByProjectID(ctx, id)
}

func (p *projectUsecase) StoreColumn(ctx context.Context, cm *domain.Column) error {
	_, err := p.projectRepo.GetByID(ctx, cm.ProjectID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}
	columns, err := p.columnRepo.FetchByProjectID(ctx, cm.ProjectID)
	if err != nil {
		return fmt.Errorf("fetch by project id: %w", err)
	}
	var ok bool
	if ok, err = isUnique(columns, cm); !ok {
		return err
	}
	if cm.Position >= len(columns) {
		cm.Position = len(columns)

		return p.columnRepo.Store(ctx, cm)
	}

	uc := columns[cm.Position:]
	for i := range uc {
		uc[i].Position++
	}
	if err = p.columnRepo.Update(ctx, uc...); err != nil {
		return fmt.Errorf("update positions: %w", err)
	}

	return p.columnRepo.Store(ctx, cm)
}

func (p *projectUsecase) FetchTasks(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
	if _, err := p.projectRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fetch tasks by project id: %w", err)
	}

	return p.taskRepo.FetchByProjectID(ctx, id)
}

func (p *projectUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	return p.projectRepo.GetByID(ctx, id)
}

func (p *projectUsecase) Update(ctx context.Context, pr *domain.Project) error {
	if _, err := p.projectRepo.GetByID(ctx, pr.ID); err != nil {
		return fmt.Errorf("update project : %w", err)
	}

	return p.projectRepo.Update(ctx, pr)
}

func (p *projectUsecase) Store(ctx context.Context, m *domain.Project) error {
	err := p.projectRepo.Store(ctx, m)
	if err != nil {
		return fmt.Errorf("store project: %w", err)
	}
	// nolint:exhaustivestruct
	err = p.columnRepo.Store(ctx, &domain.Column{
		Name:      "Default",
		Status:    "Default",
		Position:  0,
		ProjectID: m.ID,
	})
	if err != nil {
		return fmt.Errorf("store default column: %w", err)
	}

	return nil
}

func (p *projectUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := p.projectRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("delete project : %w", err)
	}

	return p.projectRepo.Delete(ctx, id)
}

func isUnique(columns []domain.Column, column *domain.Column) (bool, error) {
	for _, c := range columns {
		if column.ID == c.ID {
			continue
		}
		if column.Name == c.Name {
			return false, fmt.Errorf("column name: %w", domain.ErrUnique)
		}
		if column.Status == c.Status {
			return false, fmt.Errorf("column status: %w", domain.ErrUnique)
		}
	}

	return true, nil
}
