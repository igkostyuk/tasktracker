package usecase_test

import (
	"context"
	"errors"
	"fmt"
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

// nolint:exhaustivestruct
func TestStoreColumn(t *testing.T) {
	is := helper.New(t)

	tt := []struct {
		name string
		cl   domain.Column
	}{
		{"with position change", domain.Column{Name: "test", Position: 1}},
		{"with position change", domain.Column{Name: "test", Position: 0}},
		{"without position change", domain.Column{Name: "test", Position: 2}},
		{"with position more then number existing", domain.Column{Name: "test", Position: 3}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			columns := []domain.Column{{Name: "0", Position: 0}, {Name: "1", Position: 1}}
			mp := &mocks.ProjectRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
					return domain.Project{}, nil
				},
			}
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					return columns, nil
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return nil
				},
				StoreFunc: func(ctx context.Context, c *domain.Column) error {
					return nil
				},
			}
			u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
			err := u.StoreColumn(context.TODO(), &tc.cl)
			is.NoErr(err)
			cg := mp.GetByIDCalls()
			cf := mc.FetchByProjectIDCalls()
			cu := mc.UpdateCalls()
			cs := mc.StoreCalls()
			is.Equal(len(cg), 1)
			is.Equal(cg[0].ID, tc.cl.ID)
			is.Equal(len(cf), 1)
			is.Equal(len(cs), 1)
			is.Equal(cs[0].C, &tc.cl)
			if tc.cl.Position >= len(columns) {
				is.Equal(len(cu), 0)
			}
			if tc.cl.Position == 1 {
				is.Equal(len(cu), 1)
				is.Equal(cu[0].Cls, []domain.Column{{Name: "1", Position: 2}})
			}
			if tc.cl.Position == 0 {
				is.Equal(len(cu), 1)
				is.Equal(cu[0].Cls, []domain.Column{{Name: "0", Position: 1}, {Name: "1", Position: 2}})
			}
		})
	}
}

// nolint:exhaustivestruct
func TestStoreColumnErrors(t *testing.T) {
	is := helper.New(t)

	tt := []struct {
		name         string
		cl           domain.Column
		getByIDError error
		fetchError   error
		updateError  error
	}{
		{
			"get project by id error",
			domain.Column{Name: "nottest", Position: 1},
			fmt.Errorf("error"), nil, nil,
		},
		{
			"fetch by project id error",
			domain.Column{Name: "nottest", Position: 1},
			nil, fmt.Errorf("error"), nil,
		},
		{
			"name not unique",
			domain.Column{Name: "test", Position: 1},
			nil, nil, nil,
		},
		{
			"status not unique",
			domain.Column{Name: "nottest", Status: "testStatus", Position: 1},
			nil, nil, nil,
		},
		{
			"update error",
			domain.Column{Name: "nottest", Position: 0},
			nil, nil, fmt.Errorf("error"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			columns := []domain.Column{{ID: uuid.New(), Name: "test", Status: "testStatus", Position: 0}}
			// nolint:exhaustivestruct
			mp := &mocks.ProjectRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
					return domain.Project{}, tc.getByIDError
				},
			}
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					return columns, tc.fetchError
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return tc.updateError
				},
				StoreFunc: func(ctx context.Context, c *domain.Column) error {
					return nil
				},
			}
			u := projectUsecase.New(mp, mc, &mocks.TaskRepositoryMock{})
			err := u.StoreColumn(context.TODO(), &tc.cl)
			is.True(err != nil)
		})
	}
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
	id := uuid.New()
	project := domain.Project{ID: id, Name: "1", Description: "testDescription"}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, nil
		},
		UpdateFunc: func(ctx context.Context, pr *domain.Project) error {
			return nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	err := u.Update(context.TODO(), &project)
	is.NoErr(err)
	cg := mp.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cu := mp.UpdateCalls()
	is.Equal(len(cu), 1)
}

func TestUpdateError(t *testing.T) {
	is := helper.New(t)
	id := uuid.New()
	project := domain.Project{ID: id, Name: "1", Description: "testDescription"}
	// nolint:exhaustivestruct
	mp := &mocks.ProjectRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			return domain.Project{}, errors.New("some error")
		},
		UpdateFunc: func(ctx context.Context, pr *domain.Project) error {
			return nil
		},
	}
	u := projectUsecase.New(mp, &mocks.ColumnRepositoryMock{}, &mocks.TaskRepositoryMock{})
	err := u.Update(context.TODO(), &project)
	is.True(err != nil)
	cg := mp.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cu := mp.UpdateCalls()
	is.Equal(len(cu), 0)
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
