package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type taskUsecase struct {
	taskRepo    domain.TaskRepository
	commentRepo domain.CommentRepository
}

// New will create new a TaskUsecase object representation of domain.TaskUsecase interface.
func New(t domain.TaskRepository, c domain.CommentRepository) domain.TaskUsecase {
	return &taskUsecase{taskRepo: t, commentRepo: c}
}

func (t *taskUsecase) Fetch(ctx context.Context) ([]domain.Task, error) {
	return t.taskRepo.Fetch(ctx)
}

func (t *taskUsecase) FetchComments(ctx context.Context, id uuid.UUID) ([]domain.Comment, error) {
	if _, err := t.taskRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("fetch comments by task id: %w", err)
	}

	return t.commentRepo.FetchByTaskID(ctx, id)
}

func (t *taskUsecase) FetchByProjectID(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
	return t.taskRepo.FetchByProjectID(ctx, id)
}

func (t *taskUsecase) FetchByColumnID(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
	return t.taskRepo.FetchByProjectID(ctx, id)
}

func (t *taskUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	return t.taskRepo.GetByID(ctx, id)
}

func (t *taskUsecase) Update(ctx context.Context, pr *domain.Task) error {
	return t.taskRepo.Update(ctx, pr)
}

func (t *taskUsecase) Store(ctx context.Context, m *domain.Task) error {
	return t.taskRepo.Store(ctx, m)
}

func (t *taskUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := t.taskRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return t.taskRepo.Delete(ctx, id)
}
