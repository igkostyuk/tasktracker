package usecase

import (
	"context"

	"github.com/igkostyuk/tasktracker/domain"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

func New(c domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{commentRepo: c}
}

func (c *commentUsecase) Fetch(ctx context.Context) ([]domain.Comment, error) {
	return c.commentRepo.Fetch(ctx)
}

func (c *commentUsecase) FetchByTaskID(ctx context.Context, id string) ([]domain.Comment, error) {
	return c.commentRepo.FetchByTaskID(ctx, id)
}

func (c *commentUsecase) GetByID(ctx context.Context, id string) (domain.Comment, error) {
	return c.commentRepo.GetByID(ctx, id)
}

func (c *commentUsecase) Update(ctx context.Context, cm *domain.Comment) error {
	return c.commentRepo.Update(ctx, cm)
}

func (c *commentUsecase) Store(ctx context.Context, ct *domain.Comment) error {
	return c.commentRepo.Store(ctx, ct)
}

func (c *commentUsecase) Delete(ctx context.Context, id string) error {
	return c.commentRepo.Delete(ctx, id)
}