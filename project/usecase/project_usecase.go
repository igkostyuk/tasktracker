package usecase

import (
	"context"
	"fmt"

	"github.com/igkostyuk/tasktracker/domain"
)

type projectUsecase struct {
	projectRepo domain.ProjectRepository
	columnRepo  domain.ColumnRepository
}

func New(p domain.ProjectRepository, c domain.ColumnRepository) domain.ProjectUsecase {
	return &projectUsecase{projectRepo: p, columnRepo: c}
}

func (p *projectUsecase) Fetch(ctx context.Context) ([]domain.Project, error) {
	return p.projectRepo.Fetch(ctx)
}

func (p *projectUsecase) FetchColumns(ctx context.Context, id string) ([]domain.Column, error) {
	return p.columnRepo.FetchByProjectID(ctx, id)
}

func (p *projectUsecase) GetByID(ctx context.Context, id string) (domain.Project, error) {
	return p.projectRepo.GetByID(ctx, id)
}

func (p *projectUsecase) Update(ctx context.Context, pr *domain.Project) error {
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
		ProjectID: m.ID,
	})
	if err != nil {
		return fmt.Errorf("store default column: %w", err)
	}

	return nil
}

func (p *projectUsecase) Delete(ctx context.Context, id string) error {
	return p.projectRepo.Delete(ctx, id)
}
