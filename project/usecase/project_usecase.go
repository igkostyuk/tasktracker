package usecase

import (
	"context"

	"github.com/igkostyuk/tasktracker/domain"
)

type projectUsecase struct {
	projectRepo domain.ProjectRepository
}

func New(p domain.ProjectRepository) domain.ProjectUsecase {
	return &projectUsecase{projectRepo: p}
}

func (p *projectUsecase) Fetch(ctx context.Context) ([]domain.Project, error) {
	return p.projectRepo.Fetch(ctx)
}

func (p *projectUsecase) GetByID(ctx context.Context, id string) (domain.Project, error) {
	return p.projectRepo.GetByID(ctx, id)
}

func (p *projectUsecase) Update(ctx context.Context, pr *domain.Project) error {
	return p.projectRepo.Update(ctx, pr)
}

func (p *projectUsecase) Store(ctx context.Context, m *domain.Project) error {
	return p.projectRepo.Store(ctx, m)
}

func (p *projectUsecase) Delete(ctx context.Context, id string) error {
	return p.projectRepo.Delete(ctx, id)
}