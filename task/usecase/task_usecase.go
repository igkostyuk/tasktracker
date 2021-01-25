package usecase

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type taskUsecase struct {
	columnRepo  domain.ColumnRepository
	taskRepo    domain.TaskRepository
	commentRepo domain.CommentRepository
}

// New will create new a TaskUsecase object representation of domain.TaskUsecase interface.
func New(cl domain.ColumnRepository, t domain.TaskRepository, c domain.CommentRepository) domain.TaskUsecase {
	return &taskUsecase{columnRepo: cl, taskRepo: t, commentRepo: c}
}

func (t *taskUsecase) Fetch(ctx context.Context) ([]domain.Task, error) {
	return t.taskRepo.Fetch(ctx)
}

func (t *taskUsecase) FetchComments(ctx context.Context, id uuid.UUID) ([]domain.Comment, error) {
	if _, err := t.taskRepo.GetByID(ctx, id); err != nil {
		return nil, fmt.Errorf("get comments by task id: %w", err)
	}

	return t.commentRepo.FetchByTaskID(ctx, id)
}

func (t *taskUsecase) StoreComment(ctx context.Context, cm *domain.Comment) error {
	if _, err := t.taskRepo.GetByID(ctx, cm.TaskID); err != nil {
		return fmt.Errorf("get comments by task id: %w", err)
	}
	cm.CreatedAt = time.Now().UTC()

	return t.commentRepo.Store(ctx, cm)
}

func (t *taskUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	return t.taskRepo.GetByID(ctx, id)
}

func (t *taskUsecase) Update(ctx context.Context, ts *domain.Task) error {
	old, err := t.taskRepo.GetByID(ctx, ts.ID)
	if err != nil {
		return fmt.Errorf("fetch task by id: %w", err)
	}
	if reflect.DeepEqual(old, *ts) {
		return nil
	}
	if ts.ColumnID != old.ColumnID {
		return t.ChangeColumn(ctx, &old, ts)
	}
	tasks, err := t.taskRepo.FetchByColumnID(ctx, ts.ColumnID)
	if err != nil {
		return fmt.Errorf("fetch task by column id: %w", err)
	}

	if ts.Position > len(tasks)-1 {
		ts.Position = len(tasks) - 1
		if reflect.DeepEqual(old, *ts) {
			return nil
		}
	}
	if old.Position < ts.Position {
		return t.MoveRight(ctx, &old, ts, tasks)
	}
	if old.Position > ts.Position {
		return t.MoveLeft(ctx, &old, ts, tasks)
	}

	return t.taskRepo.Update(ctx, *ts)
}

func (t *taskUsecase) ChangeColumn(ctx context.Context, old, tk *domain.Task) error {
	_, err := t.columnRepo.GetByID(ctx, tk.ColumnID)
	if err != nil {
		return fmt.Errorf("get column by id: %w", err)
	}
	oldTasks, err := t.taskRepo.FetchByColumnID(ctx, old.ColumnID)
	if err != nil {
		return fmt.Errorf("fetch old task by column id: %w", err)
	}
	tasks, err := t.taskRepo.FetchByColumnID(ctx, tk.ColumnID)
	if err != nil {
		return fmt.Errorf("fetch task by column id: %w", err)
	}
	if tk.Position >= len(tasks) {
		tk.Position = len(tasks)
	}
	if len(oldTasks) > 0 {
		oldTasks = oldTasks[old.Position+1:]
	}
	for i := range oldTasks {
		oldTasks[i].Position--
	}
	tasks = tasks[tk.Position:]
	for i := range tasks {
		tasks[i].Position++
	}
	tasks = append(tasks, oldTasks...)
	tasks = append(tasks, *tk)

	return t.taskRepo.Update(ctx, tasks...)
}

func (t *taskUsecase) MoveRight(ctx context.Context, old, tk *domain.Task, tks []domain.Task) error {
	tks = tks[old.Position+1 : tk.Position+1]
	for i := range tks {
		tks[i].Position--
	}
	tks = append(tks, *tk)

	return t.taskRepo.Update(ctx, tks...)
}

func (t *taskUsecase) MoveLeft(ctx context.Context, old, tk *domain.Task, tks []domain.Task) error {
	tks = tks[tk.Position:old.Position]
	for i := range tks {
		tks[i].Position++
	}
	tks = append(tks, *tk)

	return t.taskRepo.Update(ctx, tks...)
}

func (t *taskUsecase) Store(ctx context.Context, tk *domain.Task) error {
	_, err := t.columnRepo.GetByID(ctx, tk.ColumnID)
	if err != nil {
		return fmt.Errorf("column get by id: %w", err)
	}
	tasks, err := t.taskRepo.FetchByColumnID(ctx, tk.ColumnID)
	if err != nil {
		return fmt.Errorf("fetch by project id: %w", err)
	}
	if tk.Position >= len(tasks) {
		tk.Position = len(tasks)

		return t.taskRepo.Store(ctx, tk)
	}

	ut := tasks[tk.Position:]
	for i := range ut {
		ut[i].Position++
	}
	if err = t.taskRepo.Update(ctx, ut...); err != nil {
		return fmt.Errorf("update positions: %w", err)
	}

	return t.taskRepo.Store(ctx, tk)
}

func (t *taskUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := t.taskRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return t.taskRepo.Delete(ctx, id)
}
