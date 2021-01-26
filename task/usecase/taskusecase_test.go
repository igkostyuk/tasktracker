package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
	mocks "github.com/igkostyuk/tasktracker/domain/mock"
	taskUsecase "github.com/igkostyuk/tasktracker/task/usecase"
	helper "github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := helper.New(t)

	want := []domain.Task{{Name: "test"}}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchFunc: func(ctx context.Context) ([]domain.Task, error) {
			return want, nil
		},
	}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
	projects, err := u.Fetch(context.TODO())
	is.NoErr(err)
	is.Equal(want, projects)
	is.Equal(len(mt.FetchCalls()), 1)
}

func TestFetchComments(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	want := []domain.Comment{{Text: "test"}}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, nil
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.CommentRepositoryMock{
		FetchByTaskIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Comment, error) {
			return want, nil
		},
	}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, mc)
	columns, err := u.FetchComments(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, columns)
	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cc := mc.FetchByTaskIDCalls()
	is.Equal(len(cc), 1)
	is.Equal(cc[0].ID, id)
}

func TestFetchTaskError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, fmt.Errorf("some error")
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.CommentRepositoryMock{
		FetchByTaskIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Comment, error) {
			return nil, nil
		},
	}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, mc)
	_, err := u.FetchComments(context.TODO(), id)
	is.True(err != nil)

	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cf := mc.FetchByTaskIDCalls()
	is.Equal(len(cf), 0)
}

func TestStoreComment(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, nil
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.CommentRepositoryMock{
		StoreFunc: func(ctx context.Context, ct *domain.Comment) error {
			return nil
		},
	}
	// nolint:exhaustivestruct
	comment := domain.Comment{TaskID: uuid.New(), Text: "test"}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, mc)
	err := u.StoreComment(context.TODO(), &comment)
	is.NoErr(err)

	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, comment.TaskID)
	cs := mc.StoreCalls()
	is.Equal(len(cs), 1)
}

func TestStoreCommentError(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, fmt.Errorf("some error")
		},
	}
	// nolint:exhaustivestruct
	mc := &mocks.CommentRepositoryMock{
		StoreFunc: func(ctx context.Context, ct *domain.Comment) error {
			return nil
		},
	}
	// nolint:exhaustivestruct
	comment := domain.Comment{TaskID: uuid.New(), Text: "test"}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, mc)
	err := u.StoreComment(context.TODO(), &comment)
	is.True(err != nil)

	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, comment.TaskID)
	cs := mc.StoreCalls()
	is.Equal(len(cs), 0)
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	want := domain.Task{Name: "test"}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return want, nil
		},
	}
	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
	project, err := u.GetByID(context.TODO(), id)
	is.NoErr(err)
	is.Equal(want, project)
	cp := mt.GetByIDCalls()
	is.Equal(len(cp), 1)
	is.Equal(cp[0].ID, id)
}

func TestMoveRight(t *testing.T) {
	tt := []struct {
		name string
		from int
		to   int
		want []domain.Task
	}{
		{
			name: "0 to 1", from: 0, to: 1,
			want: []domain.Task{
				{Name: "1", Position: 0},
				{Name: "test", Position: 1},
			},
		}, {
			name: "1 to 3", from: 1, to: 3,
			want: []domain.Task{
				{Name: "2", Position: 1},
				{Name: "3", Position: 2},
				{Name: "test", Position: 3},
			},
		}, {
			name: "0 to 5", from: 0, to: 5,
			want: []domain.Task{
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
			tasks := []domain.Task{
				{Name: "0", Position: 0},
				{Name: "1", Position: 1},
				{Name: "2", Position: 2},
				{Name: "3", Position: 3},
				{Name: "4", Position: 4},
				{Name: "5", Position: 5},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			tk := domain.Task{Name: "test", Position: tc.to}
			u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
			err := u.MoveRight(context.TODO(), &tasks[tc.from], &tk, tasks)
			is.NoErr(err)
			cu := mt.UpdateCalls()
			is.Equal(len(cu), 1)
			is.Equal(cu[0].Tks, tc.want)
		})
	}
}

// nolint:exhaustivestruct
func TestMoveLeft(t *testing.T) {
	tt := []struct {
		name string
		from int
		to   int
		want []domain.Task
	}{
		{
			name: "1 to 0", from: 1, to: 0,
			want: []domain.Task{
				{Name: "0", Position: 1},
				{Name: "test", Position: 0},
			},
		}, {
			name: "3 to 1", from: 3, to: 1,
			want: []domain.Task{
				{Name: "1", Position: 2},
				{Name: "2", Position: 3},
				{Name: "test", Position: 1},
			},
		}, {
			name: "5 to 0", from: 5, to: 0,
			want: []domain.Task{
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
			tasks := []domain.Task{
				{Name: "0", Position: 0},
				{Name: "1", Position: 1},
				{Name: "2", Position: 2},
				{Name: "3", Position: 3},
				{Name: "4", Position: 4},
				{Name: "5", Position: 5},
			}
			mt := &mocks.TaskRepositoryMock{
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
			}
			tk := domain.Task{Name: "test", Position: tc.to}
			u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
			err := u.MoveLeft(context.TODO(), &tasks[tc.from], &tk, tasks)
			is.NoErr(err)
			cu := mt.UpdateCalls()
			is.Equal(len(cu), 1)
			is.Equal(cu[0].Tks, tc.want)
		})
	}
}

func TestChangeColumn(t *testing.T) {
	is := helper.New(t)
	firstID := uuid.New()
	secondID := uuid.New()
	// nolint:exhaustivestruct
	mc := &mocks.ColumnRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
			return domain.Column{}, nil
		},
	}
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			if id == firstID {
				return []domain.Task{
					{Name: "0", Position: 0, ColumnID: firstID},
					{Name: "1", Position: 1, ColumnID: firstID},
					{Name: "2", Position: 2, ColumnID: firstID},
				}, nil
			}

			return []domain.Task{
				{Name: "0", Position: 0, ColumnID: secondID},
				{Name: "1", Position: 1, ColumnID: secondID},
				{Name: "2", Position: 2, ColumnID: secondID},
				{Name: "3", Position: 3, ColumnID: secondID},
			}, nil
		},
		UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
			return nil
		},
	}
	want := []domain.Task{
		{Name: "2", Position: 3, ColumnID: secondID},
		{Name: "3", Position: 4, ColumnID: secondID},
		{Name: "2", Position: 1, ColumnID: firstID},
		{Name: "test", Position: 2, ColumnID: secondID},
	}
	// nolint:exhaustivestruct
	otk := domain.Task{Name: "test", Position: 1, ColumnID: firstID}
	// nolint:exhaustivestruct
	tk := domain.Task{Name: "test", Position: 2, ColumnID: secondID}
	u := taskUsecase.New(mc, mt, &mocks.CommentRepositoryMock{})
	err := u.ChangeColumn(context.TODO(), &otk, &tk)
	is.NoErr(err)
	cg := mc.GetByIDCalls()
	cf := mt.FetchByColumnIDCalls()
	cu := mt.UpdateCalls()
	is.Equal(len(cg), 1)
	is.Equal(len(cf), 2)
	is.Equal(cf[0].ID, firstID)
	is.Equal(cf[1].ID, secondID)
	is.Equal(len(cu), 1)
	is.Equal(cu[0].Tks, want)
}

func TestChangeColumnError(t *testing.T) {
	tt := []struct {
		name           string
		getErr         error
		firstFetchErr  error
		secondFetchErr error
	}{
		{"get column error", fmt.Errorf("some error"), nil, nil},
		{"first fetch error", nil, fmt.Errorf("some error"), nil},
		{"second fetch error", nil, nil, fmt.Errorf("some error")},
	}
	is := helper.New(t)
	firstID := uuid.New()
	secondID := uuid.New()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return domain.Column{}, tc.getErr
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					if id == firstID {
						return []domain.Task{}, tc.firstFetchErr
					}

					return []domain.Task{}, tc.secondFetchErr
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
			}

			// nolint:exhaustivestruct
			otk := domain.Task{Name: "test", ColumnID: firstID}
			// nolint:exhaustivestruct
			tk := domain.Task{Name: "test", ColumnID: secondID}
			u := taskUsecase.New(mc, mt, &mocks.CommentRepositoryMock{})
			err := u.ChangeColumn(context.TODO(), &otk, &tk)
			is.True(err != nil)
		})
	}
}

//nolint:exhaustivestruct,funlen
func TestUpdate(t *testing.T) {
	tt := []struct {
		name string
		old  domain.Task
		tk   domain.Task
	}{
		{"tasks equal", domain.Task{}, domain.Task{}},
		{
			"position more then columns len and rest equal",
			domain.Task{Position: 2},
			domain.Task{Position: 3},
		},
		{
			"different column ID and position more then len",
			domain.Task{Position: 1},
			domain.Task{ColumnID: uuid.New(), Position: 3},
		},
		{"move right", domain.Task{Position: 1}, domain.Task{Position: 3}},
		{"move left", domain.Task{Position: 1}, domain.Task{Position: 0}},
		{"update name", domain.Task{Name: "nottest"}, domain.Task{Name: "test"}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return domain.Column{}, nil
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
					return tc.old, nil
				},
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					if id == tc.old.ColumnID {
						return []domain.Task{{Name: "0", Position: 0}, tc.old, {Name: "2", Position: 2}}, nil
					}

					return []domain.Task{
						{Name: "0", Position: 0, ColumnID: tc.tk.ColumnID},
						{Name: "1", Position: 1, ColumnID: tc.tk.ColumnID},
						{Name: "2", Position: 2, ColumnID: tc.tk.ColumnID},
					}, nil
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
			}
			u := taskUsecase.New(mc, mt, &mocks.CommentRepositoryMock{})
			err := u.Update(context.TODO(), &tc.tk)
			is.NoErr(err)
			cg := mt.GetByIDCalls()
			cf := mt.FetchByColumnIDCalls()
			cu := mt.UpdateCalls()
			is.Equal(len(cg), 1)
			is.Equal(cg[0].ID, tc.tk.ID)
			if reflect.DeepEqual(tc.old, tc.tk) {
				is.Equal(len(cu), 0)
				if tc.old.Position >= 2 {
					is.Equal(len(cf), 1)
					is.Equal(cf[0].ID, tc.tk.ColumnID)

					return
				}
				is.Equal(len(cf), 0)

				return
			}
			if tc.old.ColumnID != tc.tk.ColumnID {
				is.Equal(len(cf), 2)
				is.Equal(cf[0].ID, tc.old.ColumnID)
				is.Equal(cf[1].ID, tc.tk.ColumnID)
				is.Equal(len(cu), 1)
				is.Equal(cu[0].Tks, []domain.Task{{Name: "2", Position: 1}, tc.tk})

				return
			}
			is.Equal(len(cf), 1)
			is.Equal(cf[0].ID, tc.tk.ColumnID)
			is.Equal(len(cu), 1)
			if tc.old.Position > tc.tk.Position {
				is.Equal(cu[0].Tks, []domain.Task{{Name: "0", Position: 1}, tc.tk})
			}
			if tc.old.Position < tc.tk.Position {
				is.Equal(cu[0].Tks, []domain.Task{{Name: "2", Position: 1}, tc.tk})
			}
			if tc.old.Name != tc.tk.Name {
				is.Equal(cu[0].Tks, []domain.Task{tc.tk})
			}
		})
	}
}

func TestUpdateError(t *testing.T) {
	tt := []struct {
		name       string
		taskName   string
		getError   error
		fetchError error
	}{
		{"get by id error", "nottest", fmt.Errorf("some error"), nil},
		{"fetch by column id error", "nottest", nil, fmt.Errorf("some error")},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
					return domain.Task{}, tc.getError
				},
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					return []domain.Task{{Name: "test"}}, tc.fetchError
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
			}
			// nolint:exhaustivestruct
			task := domain.Task{Name: tc.taskName, ID: uuid.New()}
			u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
			err := u.Update(context.TODO(), &task)
			is.True(err != nil)
		})
	}
}

func TestStore(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	tt := []struct {
		name string
		tk   domain.Task
	}{
		{"with position change", domain.Task{Name: "test", Position: 1}},
		{"with position change", domain.Task{Name: "test", Position: 0}},
		{"without position change", domain.Task{Name: "test", Position: 2}},
		{"with position more then number existing", domain.Task{Name: "test", Position: 3}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tasks := []domain.Task{{Name: "0", Position: 0}, {Name: "1", Position: 1}}
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return domain.Column{}, nil
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					return tasks, nil
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return nil
				},
				StoreFunc: func(ctx context.Context, c *domain.Task) error {
					return nil
				},
			}
			u := taskUsecase.New(mc, mt, &mocks.CommentRepositoryMock{})
			err := u.Store(context.TODO(), &tc.tk)
			is.NoErr(err)
			cg := mc.GetByIDCalls()
			cf := mt.FetchByColumnIDCalls()
			cu := mt.UpdateCalls()
			cs := mt.StoreCalls()
			is.Equal(len(cg), 1)
			is.Equal(cg[0].ID, tc.tk.ID)
			is.Equal(len(cf), 1)
			is.Equal(len(cs), 1)
			is.Equal(cs[0].T, &tc.tk)
			if tc.tk.Position >= len(tasks) {
				is.Equal(len(cu), 0)
			}
			if tc.tk.Position == 1 {
				is.Equal(len(cu), 1)
				is.Equal(cu[0].Tks, []domain.Task{{Name: "1", Position: 2}})
			}
			if tc.tk.Position == 0 {
				is.Equal(len(cu), 1)
				is.Equal(cu[0].Tks, []domain.Task{{Name: "0", Position: 1}, {Name: "1", Position: 2}})
			}
		})
	}
}

//nolint:exhaustivestruct,funlen
func TestStoreErrors(t *testing.T) {
	is := helper.New(t)

	tt := []struct {
		name         string
		tk           domain.Task
		getByIDError error
		fetchError   error
		updateError  error
	}{
		{
			"get project by id error",
			domain.Task{Name: "nottest", Position: 1},
			fmt.Errorf("error"), nil, nil,
		},
		{
			"fetch by project id error",
			domain.Task{Name: "nottest", Position: 1},
			nil, fmt.Errorf("error"), nil,
		},
		{
			"update error",
			domain.Task{Name: "nottest", Position: 0},
			nil, nil, fmt.Errorf("error"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tasks := []domain.Task{{ID: uuid.New(), Name: "test", Position: 0}}
			// nolint:exhaustivestruct
			mc := &mocks.ColumnRepositoryMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Column, error) {
					return domain.Column{}, tc.getByIDError
				},
			}
			// nolint:exhaustivestruct
			mt := &mocks.TaskRepositoryMock{
				FetchByColumnIDFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					return tasks, tc.fetchError
				},
				UpdateFunc: func(ctx context.Context, cls ...domain.Task) error {
					return tc.updateError
				},
				StoreFunc: func(ctx context.Context, c *domain.Task) error {
					return nil
				},
			}
			u := taskUsecase.New(mc, mt, &mocks.CommentRepositoryMock{})
			err := u.Store(context.TODO(), &tc.tk)
			is.True(err != nil)
		})
	}
}

func TestDelete(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, nil
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
	err := u.Delete(context.TODO(), id)
	is.NoErr(err)

	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cd := mt.DeleteCalls()
	is.Equal(len(cd), 1)
	is.Equal(cd[0].ID, id)
}

func TestDeleteError(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	// nolint:exhaustivestruct
	mt := &mocks.TaskRepositoryMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Task, error) {
			return domain.Task{}, errors.New("some error")
		},
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	u := taskUsecase.New(&mocks.ColumnRepositoryMock{}, mt, &mocks.CommentRepositoryMock{})
	err := u.Delete(context.TODO(), id)
	is.True(err != nil)

	cg := mt.GetByIDCalls()
	is.Equal(len(cg), 1)
	is.Equal(cg[0].ID, id)
	cd := mt.DeleteCalls()
	is.Equal(len(cd), 0)
}
