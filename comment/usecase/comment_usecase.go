package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
}

// New will create new a CommentUsecase object representation of domain.ComentUsecase interface.
func New(c domain.CommentRepository) domain.CommentUsecase {
	return &commentUsecase{commentRepo: c}
}

func (c *commentUsecase) Fetch(ctx context.Context) ([]domain.Comment, error) {
	return c.commentRepo.Fetch(ctx)
}

func (c *commentUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
	return c.commentRepo.GetByID(ctx, id)
}

func (c *commentUsecase) Update(ctx context.Context, cm *domain.Comment) error {
	old, err := c.commentRepo.GetByID(ctx, cm.ID)
	if err != nil {
		return fmt.Errorf("get by id comment: %w", err)
	}
	cm.TaskID = old.TaskID
	cm.CreatedAt = old.CreatedAt
	return c.commentRepo.Update(ctx, cm)
}

func (c *commentUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := c.commentRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get by id comment: %w", err)
	}

	return c.commentRepo.Delete(ctx, id)
}
