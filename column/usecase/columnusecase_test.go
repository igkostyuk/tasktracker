package usecase_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	columnUsecase "github.com/igkostyuk/tasktracker/column/usecase"
	"github.com/igkostyuk/tasktracker/domain"
	mocks "github.com/igkostyuk/tasktracker/domain/mock"
	helper "github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := helper.New(t)

	want := []domain.Column{{Name: "test", Status: "testStatus"}}
	// nolint:exhaustivestruct
	mockedColumnRepo := &mocks.ColumnRepositoryMock{
		FetchFunc: func(ctx context.Context) ([]domain.Column, error) {
			return want, nil
		},
	}
	u := columnUsecase.New(mockedColumnRepo, &mocks.TaskRepositoryMock{})
	projects, err := u.Fetch(context.TODO())
	is.NoErr(err)
	is.Equal(want, projects)
	is.Equal(len(mockedColumnRepo.FetchCalls()), 1)
}

func TestFetchTask(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	want := []domain.Task{{Name: "1", Position: 1, Description: "testDescription"}}
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
			return domain.Column{}, nil
		},
	}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			return want, nil
		},
	}
	u := columnUsecase.New(mc, mt)
	columns, err := u.FetchTasks(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, columns)
	cp := mc.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mt.FetchByColumnIDCalls()
	is.Equal(len(cc), 1)
	is.Equal(cc[0].ID, id)
}

func TestFetchTaskError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
			return domain.Column{}, fmt.Errorf("some error")
		},
	}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			return nil, nil
		},
	}
	u := columnUsecase.New(mc, mt)
	_, err := u.FetchTasks(context.TODO(), id)
	is.True(err != nil)

	cp := mc.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
	cc := mt.FetchByColumnIDCalls()
	is.Equal(len(cc), 0)
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	want := domain.Column{Name: "test", Status: "testStatus"}
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
			return want, nil
		},
	}
	u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
	project, err := u.GetByID(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, project)
	cp := mc.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
}

func TestFetchByProjectID(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
			return []domain.Column{}, nil
		},
	}
	u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
	_, err := u.FetchByProjectID(context.TODO(), id)
	is.NoErr(err)
	cf := mc.FetchByProjectIDCalls()
	is.Equal(len(cf), 1)
	is.Equal(cf[0].ID, id)
}

func TestMoveRight(t *testing.T) {
	tt := []struct {
		name string
		from int
		to   int
		want []domain.Column
	}{
		{
			name: "0 to 1", from: 0, to: 1,
			want: []domain.Column{
				{Name: "1", Position: 0},
				{Name: "test", Position: 1},
			},
		}, {
			name: "1 to 3", from: 1, to: 3,
			want: []domain.Column{
				{Name: "2", Position: 1},
				{Name: "3", Position: 2},
				{Name: "test", Position: 3},
			},
		}, {
			name: "0 to 5", from: 0, to: 5,
			want: []domain.Column{
				{Name: "1", Position: 0},
				{Name: "2", Position: 1},
				{Name: "3", Position: 2},
				{Name: "4", Position: 3},
				{Name: "5", Position: 4},
				{Name: "test", Position: 5},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			columns := []domain.Column{
				{Name: "0", Position: 0},
				{Name: "1", Position: 1},
				{Name: "2", Position: 2},
				{Name: "3", Position: 3},
				{Name: "4", Position: 4},
				{Name: "5", Position: 5},
			}
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			cl := domain.Column{Name: "test", Position: tc.to}
			u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
			err := u.MoveRight(context.TODO(), &columns[tc.from], &cl, columns)
			is.NoErr(err)
			cu := mc.UpdateCalls()
			is.Equal(len(cu), 1)
			is.Equal(cu[0].Cls, tc.want)
		})
	}
}

// nolint:exhaustivestruct
func TestMoveLeft(t *testing.T) {
	tt := []struct {
		name string
		from int
		to   int
		want []domain.Column
	}{
		{
			name: "1 to 0", from: 1, to: 0,
			want: []domain.Column{
				{Name: "0", Position: 1},
				{Name: "test", Position: 0},
			},
		}, {
			name: "3 to 1", from: 3, to: 1,
			want: []domain.Column{
				{Name: "1", Position: 2},
				{Name: "2", Position: 3},
				{Name: "test", Position: 1},
			},
		}, {
			name: "5 to 0", from: 5, to: 0,
			want: []domain.Column{
				{Name: "0", Position: 1},
				{Name: "1", Position: 2},
				{Name: "2", Position: 3},
				{Name: "3", Position: 4},
				{Name: "4", Position: 5},
				{Name: "test", Position: 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			columns := []domain.Column{
				{Name: "0", Position: 0},
				{Name: "1", Position: 1},
				{Name: "2", Position: 2},
				{Name: "3", Position: 3},
				{Name: "4", Position: 4},
				{Name: "5", Position: 5},
			}
			mc := &mocks.ColumnRepositoryMock{
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return nil
				},
			}
			cl := domain.Column{Name: "test", Position: tc.to}
			u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
			err := u.MoveLeft(context.TODO(), &columns[tc.from], &cl, columns)
			is.NoErr(err)
			cu := mc.UpdateCalls()
			is.Equal(len(cu), 1)
			is.Equal(cu[0].Cls, tc.want)
		})
	}
}

// nolint:exhaustivestruct,funlen
func TestUpdate(t *testing.T) {
	tt := []struct {
		name string
		old  domain.Column
		cl   domain.Column
	}{
		{"equal", domain.Column{}, domain.Column{}},
		{"position more then columns len and rest equal", domain.Column{Position: 2}, domain.Column{Position: 3}},
		{"move right", domain.Column{Position: 1}, domain.Column{Position: 3}},
		{"move left", domain.Column{Position: 1}, domain.Column{Position: 0}},
		{"update name", domain.Column{Name: "nottest"}, domain.Column{Name: "test"}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return tc.old, nil
				},
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					return []domain.Column{{Name: "0", Position: 0}, tc.old, {Name: "2", Position: 2}}, nil
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return nil
				},
			}
			u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
			err := u.Update(context.TODO(), &tc.cl)
			is.NoErr(err)
			cg := mc.GetByIDCalls()
			cf := mc.FetchByProjectIDCalls()
			cu := mc.UpdateCalls()
			is.Equal(len(cg), 1)
			is.Equal(cg[0].ID, tc.cl.ID)
			if reflect.DeepEqual(tc.old, tc.cl) {
				is.Equal(len(cu), 0)
				if tc.old.Position >= 2 {
					is.Equal(len(cf), 1)
					is.Equal(cf[0].ID, tc.cl.ProjectID)

					return
				}
				is.Equal(len(cf), 0)

				return
			}
			is.Equal(len(cf), 1)
			is.Equal(cf[0].ID, tc.cl.ProjectID)
			is.Equal(len(cu), 1)
			if tc.old.Position > tc.cl.Position {
				is.Equal(cu[0].Cls, []domain.Column{{Name: "0", Position: 1}, tc.cl})
			}
			if tc.old.Position < tc.cl.Position {
				is.Equal(cu[0].Cls, []domain.Column{{Name: "2", Position: 1}, tc.cl})
			}
			if tc.old.Name != tc.cl.Name {
				is.Equal(cu[0].Cls, []domain.Column{tc.cl})
			}
		})
	}
}

func TestUpdateError(t *testing.T) {
	tt := []struct {
		name       string
		columnName string
		getError   error
		fetchError error
	}{
		{"get by id error", "nottest", fmt.Errorf("some error"), nil},
		{"fetch by project id error", "nottest", nil, fmt.Errorf("some error")},
		{"name not unique", "test", nil, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return domain.Column{}, tc.getError
				},
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					return []domain.Column{{Name: "test"}}, tc.fetchError
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Column) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			column := domain.Column{Name: tc.columnName, ID: uuid.New()}
			u := columnUsecase.New(mc, &mocks.TaskRepositoryMock{})
			err := u.Update(context.TODO(), &column)
			is.True(err != nil)
		})
	}
}

// nolint:funlen
func TestDelete(t *testing.T) {
	is := helper.New(t)
	firstColumnID := uuid.New()
	secondColumnID := uuid.New()
	projectID := uuid.New()

	tt := []struct {
		name     string
		columnID uuid.UUID
	}{
		{"without tasks", firstColumnID},
		{"with tasks", secondColumnID},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					if id == firstColumnID {
						// nolint:exhaustivestruct
						return domain.Column{ID: tc.columnID, ProjectID: projectID, Position: 0}, nil
					}

					return domain.Column{ID: tc.columnID, ProjectID: projectID, Position: 1}, nil
				},
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					return []domain.Column{{ID: firstColumnID}, {ID: secondColumnID}}, nil
				},
				DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					if id == firstColumnID {
						return nil, nil
					}

					return []domain.Task{{ColumnID: secondColumnID}}, nil
				},
				UpdateFunc: func(ctx context.Context, tks ...domain.Task) error {
					return nil
				},
			}
			u := columnUsecase.New(mc, mt)
			err := u.Delete(context.TODO(), tc.columnID)
			is.NoErr(err)
			ccg := mc.GetByIDCalls()
			ccf := mc.FetchByProjectIDCalls()
			ccd := mc.DeleteCalls()
			ctf := mt.FetchByColumnIDCalls()
			ctu := mt.UpdateCalls()
			is.Equal(len(ccg), 1)
			is.Equal(ccg[0].ID, tc.columnID)
			is.Equal(len(ccf), 1)
			is.Equal(ccf[0].ID, projectID)
			is.Equal(len(ccd), 1)
			is.Equal(ccd[0].ID, tc.columnID)
			is.Equal(ctf[0].ID, tc.columnID)
			if tc.columnID == firstColumnID {
				is.Equal(len(ctf), 1)
				is.Equal(len(ctu), 0)
			}
			if tc.columnID == secondColumnID {
				is.Equal(len(ctf), 2)
				is.Equal(ctf[1].ID, firstColumnID)
				is.Equal(len(ctu), 1)
				is.Equal(ctu[0].Tks, []domain.Task{{ColumnID: firstColumnID}})
			}
		})
	}
}

// nolint:funlen
func TestDeleteError(t *testing.T) {
	is := helper.New(t)
	firstColumnID := uuid.New()
	secondColumnID := uuid.New()
	thirdColumnID := uuid.New()
	projectID := uuid.New()

	tt := []struct {
		name                        string
		columnID                    uuid.UUID
		columnGetByIDError          error
		columnFetchByProjectIDError error
		taskFetchByColumnError      error
		taskUpdateError             error
	}{
		{"column get by id error", secondColumnID, fmt.Errorf("some error"), nil, nil, nil},
		{"column fetch by project id error", secondColumnID, nil, fmt.Errorf("some error"), nil, nil},
		{"task fetch by column id error", secondColumnID, nil, nil, fmt.Errorf("some error"), nil},
		{"task update error", secondColumnID, nil, nil, nil, fmt.Errorf("some error")},
		{"task update error", thirdColumnID, nil, nil, fmt.Errorf("some error"), nil},
		{"last column", firstColumnID, nil, nil, nil, fmt.Errorf("some error")},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					if id == firstColumnID {
						// nolint:exhaustivestruct
						return domain.Column{ID: tc.columnID, Position: 0}, tc.columnGetByIDError
					}

					return domain.Column{ID: tc.columnID, ProjectID: projectID, Position: 1}, tc.columnGetByIDError
				},
				FetchByProjectIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					if id != projectID {
						return []domain.Column{{ID: firstColumnID}}, nil
					}

					return []domain.Column{{ID: firstColumnID}, {ID: secondColumnID}}, tc.columnFetchByProjectIDError
				},
				DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					if id == firstColumnID {
						return nil, tc.taskFetchByColumnError
					}
					if id == thirdColumnID {
						return nil, tc.taskFetchByColumnError
					}

					return []domain.Task{{ColumnID: secondColumnID}}, nil
				},
				UpdateFunc: func(ctx context.Context, tks ...domain.Task) error {
					return tc.taskUpdateError
				},
			}
			u := columnUsecase.New(mc, mt)
			err := u.Delete(context.TODO(), tc.columnID)
			is.True(err != nil)
		})
	}
}
