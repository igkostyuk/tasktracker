package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
	mocks "github.com/igkostyuk/tasktracker/domain/mock"
	projectUsecase "github.com/igkostyuk/tasktracker/project/usecase"
	helper "github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := helper.New(t)

	want := []domain.Project{{Name: "1", Description: "testDescription"}}
	// nolint:exhaustivestruct
	mockedProjectRepo := &mocks.ProjectRepositoryMock{
		FetchFunc: func(ctx context.Context) ([]domain.Project, error) {
			return want, nil
		},
	}
	u := projectUsecase.New(mockedProjectRepo, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	projects, err := u.Fetch(context.TODO())
	is.NoErr(err)
	is.Equal(want, projects)
	is.Equal(len(mockedProjectRepo.FetchCalls()), 1)
}

func TestFetchColumns(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	want := []domain.Column{{Name: "1", Position: 1, Status: "testStatus"}}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, nil
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
			return want, nil
		},
	}
	u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
	projects, err := u.FetchColumns(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, projects)
	cp := mp.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mc.FetchByProjectIDCalls()
	is.Equal(len(cc), 1)
	is.Equal(cc[0].ID, id)
}

func TestFetchColumnsError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, errors.New("some error")
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
			return nil, nil
		},
	}
	u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
	_, err := u.FetchColumns(context.TODO(), id)
	is.True(err != nil)
	cp := mp.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mc.FetchByProjectIDCalls()
	is.Equal(len(cc), 0)
}

func TestFetchTask(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	want := []domain.Task{{Name: "1", Position: 1, Description: "testDescription"}}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, nil
		},
	}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			return want, nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, mt)
	projects, err := u.FetchTasks(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, projects)
	cp := mp.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mt.FetchByProjectIDCalls()
	is.Equal(len(cc), 1)
	is.Equal(cc[0].ID, id)
}

func TestFetchTaskError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, errors.New("some error")
		},
	}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			return nil, nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, mt)
	_, err := u.FetchTasks(context.TODO(), id)
	is.True(err != nil)

	cp := mp.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mt.FetchByProjectIDCalls()
	is.Equal(len(cc), 0)
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	want := domain.Project{Name: "1", Description: "testDescription"}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return want, nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	project, err := u.GetByID(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, project)
	cp := mp.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
}

func TestUpdate(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	project := domain.Project{Name: "1", Description: "testDescription"}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		UpdateFunc: func(ctx context.Context, pr *domain.Project) error {
			return nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	err := u.Update(context.TODO(), &project)
	is.NoErr(err)
	cp := mp.UpdateCalls()
	is.Equal(len(cp), 1)
}

func TestStore(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	project := domain.Project{ID: id}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		StoreFunc: func(ctx context.Context, p *domain.Project) error {
			return nil
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		StoreFunc: func(ctx context.Context, c *domain.Column) error {
			return nil
		},
	}

	u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
	err := u.Store(context.TODO(), &project)
	is.NoErr(err)

	cp := mp.StoreCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].A, &project)
	cc := mc.StoreCalls()
	is.Equal(len(cc), 1)
	is.Equal(cc[0].C.ProjectID, id)
}

//nolint:funlen
func TestStoreErrors(t *testing.T) {
	is := helper.New(t)
	id := uuid.New()
	// nolint:exhaustivestruct
	project := domain.Project{ID: id}
	t.Run("project store error", func(t *testing.T) {
		// nolint:exhaustivestruct
		mp := &mocks.ProjectRepositoryMock{
			StoreFunc: func(ctx context.Context, p *domain.Project) error {
				return errors.New("some error")
			},
		}
		// nolint:exhaustivestruct
		mc := &mocks.ColumnRepositoryMock{
			StoreFunc: func(ctx context.Context, c *domain.Column) error {
				return nil
			},
		}
		u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
		err := u.Store(context.TODO(), &project)
		is.True(err != nil)
		cp := mp.StoreCalls()
		is.Equal(len(cp), 1)
		is.Equal(cp[0].A, &project)
		cc := mc.StoreCalls()
		is.Equal(len(cc), 0)
	})
	t.Run("column store error", func(t *testing.T) {
		// nolint:exhaustivestruct
		mp := &mocks.ProjectRepositoryMock{
			StoreFunc: func(ctx context.Context, p *domain.Project) error {
				return nil
			},
		}
		// nolint:exhaustivestruct
		mc := &mocks.ColumnRepositoryMock{
			StoreFunc: func(ctx context.Context, c *domain.Column) error {
				return errors.New("some error")
			},
		}
		u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
		err := u.Store(context.TODO(), &project)
		is.True(err != nil)
		cp := mp.StoreCalls()
		is.Equal(len(cp), 1)
		is.Equal(cp[0].A, &project)
		cc := mc.StoreCalls()
		is.Equal(len(cc), 1)
		is.Equal(cc[0].C.ProjectID, id)
	})
}

func TestDelete(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, nil
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	err := u.Delete(context.TODO(), id)
	is.NoErr(err)

	cg := mp.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cd := mp.DeleteCalls()
	is.Equal(len(cd), 1)
	is.Equal(cd[0].ID, id)
}

func TestDeleteError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, errors.New("some error")
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	err := u.Delete(context.TODO(), id)
	is.True(err != nil)

	cg := mp.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cd := mp.DeleteCalls()
	is.Equal(len(cd), 0)
}
