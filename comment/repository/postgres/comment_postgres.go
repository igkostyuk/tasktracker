package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/igkostyuk/tasktracker/domain"
)

type commentRepository struct {
	db *sql.DB
}

// New will create new a CommentRepository object representation of domain.CommentRepository interface.
func New(db *sql.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (c *commentRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Comment, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	result := make([]domain.Comment, 0)
	for rows.Next() {
		t := domain.Comment{}
		err = rows.Scan(
			&t.ID,
			&t.Text,
			&t.TaskID,
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

func (c *commentRepository) Fetch(ctx context.Context) ([]domain.Comment, error) {
	query := `SELECT id, text, task_id FROM comments`

	return c.fetch(ctx, query)
}

func (c *commentRepository) FetchByTaskID(ctx context.Context, id string) ([]domain.Comment, error) {
	query := `SELECT id, text, task_id FROM comments WHERE task_id = $1`

	return c.fetch(ctx, query, id)
}

func (c *commentRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Comment, error) {
	row := c.db.QueryRowContext(ctx, query, args...)
	res := domain.Comment{}
	err := row.Scan(
		&res.ID,
		&res.Text,
		&res.TaskID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Comment{}, fmt.Errorf("comment: %w", domain.ErrNotFound)
	}
	if err != nil {
		return domain.Comment{}, fmt.Errorf("getOne error: %w", err)
	}

	return res, nil
}

func (c *commentRepository) GetByID(ctx context.Context, id string) (domain.Comment, error) {
	query := `SELECT id, position, name, description, column_id FROM comments WHERE id = $1`

	return c.getOne(ctx, query, id)
}

func (c *commentRepository) Update(ctx context.Context, cm *domain.Comment) error {
	query := `UPDATE comments SET taxt=$2, task_id=$3 FROM comments WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, cm.ID, cm.Text, cm.TaskID)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (c *commentRepository) Store(ctx context.Context, ct *domain.Comment) error {
	query := `INSERT INTO comments (text, task_id) VALUES ( $1, $2) RETURNING id`
	row := c.db.QueryRowContext(ctx, query, ct.Text, ct.TaskID)
	err := row.Scan(&ct.ID)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	return nil
}

func (c *commentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}
