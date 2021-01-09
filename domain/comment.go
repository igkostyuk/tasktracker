package domain

import (
	"context"
)

//go:generate moq -out ./mock/comment.go -pkg mocks . CommentUsecase CommentRepository

// Comment represent a comment in tasktracker.
type Comment struct {
	ID     string `json:"id" readonly:"true"`
	Text   string `json:"text" validate:"required,min=1,max=5000"`
	TaskID string `json:"task_id" validate:"required,uuid4"`
}

// CommentUsecase represent the comment's usecases.
type CommentUsecase interface {
	Fetch(ctx context.Context) ([]Comment, error)
	FetchByTaskID(ctx context.Context, id string) ([]Comment, error)
	GetByID(ctx context.Context, id string) (Comment, error)
	Update(ctx context.Context, tk *Comment) error
	Store(context.Context, *Comment) error
	Delete(ctx context.Context, id string) error
}

// CommentRepository represent the comment's repository contract.
type CommentRepository interface {
	Fetch(ctx context.Context) ([]Comment, error)
	FetchByTaskID(ctx context.Context, id string) ([]Comment, error)
	GetByID(ctx context.Context, id string) (Comment, error)
	Update(ctx context.Context, cm *Comment) error
	Store(ctx context.Context, ct *Comment) error
	Delete(ctx context.Context, id string) error
}
