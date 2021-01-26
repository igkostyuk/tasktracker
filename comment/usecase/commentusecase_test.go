package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	commentUsecase "github.com/igkostyuk/tasktracker/comment/usecase"
	"github.com/igkostyuk/tasktracker/domain"
	mocks "github.com/igkostyuk/tasktracker/domain/mock"
	helper "github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := helper.New(t)

	want := []domain.Comment{{Text: "test"}}
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		FetchFunc: func(ctx context.Context) ([]domain.Comment, error) {
			return want, nil
		},
	}
	u := commentUsecase.New(mockedCommentRepo)
	projects, err := u.Fetch(context.TODO())
	is.NoErr(err)
	is.Equal(want, projects)
	is.Equal(len(mockedCommentRepo.FetchCalls()), 1)
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)
	id := uuid.New()
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
			return domain.Comment{}, nil
		},
	}
	u := commentUsecase.New(mockedCommentRepo)
	_, err := u.GetByID(context.TODO(), id)
	is.NoErr(err)
	cg := mockedCommentRepo.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
}

func TestUpdate(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
			return domain.Comment{}, nil
		},
		UpdateFunc: func(ctx context.Context, cm *domain.Comment) error {
			return nil
		},
	}
	// nolint:exhaustivestruct
	comment := domain.Comment{ID: uuid.New(), Text: "test"}
	u := commentUsecase.New(mockedCommentRepo)
	err := u.Update(context.TODO(), &comment)
	is.NoErr(err)
	cg := mockedCommentRepo.GetByIDCalls()
	cu := mockedCommentRepo.UpdateCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, comment.ID)
	is.Equal(len(cu), 1)
}

func TestUpdateError(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
			return domain.Comment{}, fmt.Errorf("some error")
		},
		UpdateFunc: func(ctx context.Context, cm *domain.Comment) error {
			return nil
		},
	}
	// nolint:exhaustivestruct
	comment := domain.Comment{ID: uuid.New(), Text: "test"}
	u := commentUsecase.New(mockedCommentRepo)
	err := u.Update(context.TODO(), &comment)
	is.True(err != nil)
	cg := mockedCommentRepo.GetByIDCalls()
	cu := mockedCommentRepo.UpdateCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, comment.ID)
	is.Equal(len(cu), 0)
}

func TestDelete(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
			return domain.Comment{}, nil
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	id := uuid.New()
	u := commentUsecase.New(mockedCommentRepo)
	err := u.Delete(context.TODO(), id)
	is.NoErr(err)
	cg := mockedCommentRepo.GetByIDCalls()
	cu := mockedCommentRepo.DeleteCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	is.Equal(len(cu), 1)
}

func TestDeleteError(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mockedCommentRepo := &mocks.CommentRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Comment, error) {
			return domain.Comment{}, fmt.Errorf("some error")
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	id := uuid.New()
	u := commentUsecase.New(mockedCommentRepo)
	err := u.Delete(context.TODO(), id)
	is.True(err != nil)
	cg := mockedCommentRepo.GetByIDCalls()
	cu := mockedCommentRepo.DeleteCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	is.Equal(len(cu), 0)
}
