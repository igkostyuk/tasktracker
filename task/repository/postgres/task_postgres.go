package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/igkostyuk/tasktracker/domain"
)

type taskRepository struct {
	db *sql.DB
}

func New(db *sql.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Task, error) {
	rows, err := t.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	result := make([]domain.Task, 0)
	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(
			&t.ID,
			&t.Position,
			&t.Name,
			&t.Description,
			&t.ColumnID,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan error: %w", err)
		}
		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered during iteration %w", err)
	}

	return result, nil
}

func (t *taskRepository) Fetch(ctx context.Context) ([]domain.Task, error) {
	query := `SELECT id, position, name, description, column_id FROM tasks`

	return t.fetch(ctx, query)
}

func (t *taskRepository) FetchByColumnID(ctx context.Context, id string) ([]domain.Task, error) {
	query := `SELECT id, position, name, description, column_id FROM tasks WHERE column_id = $1`

	return t.fetch(ctx, query, id)
}

func (t *taskRepository) FetchByProjectID(ctx context.Context, id string) ([]domain.Task, error) {
	query := `SELECT id, position, name, description, column_id FROM tasks WHERE project_id = $1`

	return t.fetch(ctx, query, id)
}

func (t *taskRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Task, error) {
	row := t.db.QueryRowContext(ctx, query, args...)
	res := domain.Task{}
	err := row.Scan(
		&res.ID,
		&res.Position,
		&res.Name,
		&res.Description,
		&res.ColumnID,
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("getOne error: %w", err)
	}

	return res, nil
}

func (t *taskRepository) GetByID(ctx context.Context, id string) (domain.Task, error) {
	query := `SELECT id, position, name, description, column_id FROM tasks WHERE id = $1`

	return t.getOne(ctx, query, id)
}

func (t *taskRepository) Update(ctx context.Context, tk *domain.Task) error {
	query := `UPDATE tasks SET position=$2, name=$3, description=$4, column_id=$5 FROM tasks WHERE id = $1`
	_, err := t.db.ExecContext(ctx, query, tk.ID, tk.Position, tk.Name, tk.Description, tk.ColumnID)

	return fmt.Errorf("update error: %w", err)
}

func (t *taskRepository) Store(ctx context.Context, ts *domain.Task) error {
	query := `INSERT INTO tasks (position, name, description, column_id) VALUES ( $1, $2, $3, $4) RETURNING id`
	row := t.db.QueryRowContext(ctx, query, ts.Position, ts.Name, ts.Description, ts.ColumnID)
	err := row.Scan(&ts.ID)

	return fmt.Errorf("store error: %w", err)
}

func (t *taskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := t.db.ExecContext(ctx, query, id)

	return fmt.Errorf("delete error: %w", err)
}
