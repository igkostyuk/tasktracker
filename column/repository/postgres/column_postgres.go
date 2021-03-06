package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type columnRepository struct {
	db *sql.DB
}

// New will create new a ColumnRepository object representation of domain.ColumnRepository interface.
func New(db *sql.DB) domain.ColumnRepository {
	return &columnRepository{db: db}
}

func (c *columnRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Column, error) {
	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	result := make([]domain.Column, 0)
	for rows.Next() {
		t := domain.Column{}
		err = rows.Scan(
			&t.ID,
			&t.Position,
			&t.Name,
			&t.Status,
			&t.ProjectID,
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

func (c *columnRepository) Fetch(ctx context.Context) ([]domain.Column, error) {
	query := `SELECT id,position,name,status,project_id FROM columns ORDER BY position`

	return c.fetch(ctx, query)
}

func (c *columnRepository) FetchByProjectID(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
	query := `SELECT id,position,name,status,project_id FROM columns WHERE project_id = $1 ORDER BY position`

	return c.fetch(ctx, query, id)
}

func (c *columnRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Column, error) {
	row := c.db.QueryRowContext(ctx, query, args...)
	res := domain.Column{}
	err := row.Scan(
		&res.ID,
		&res.Position,
		&res.Name,
		&res.Status,
		&res.ProjectID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Column{}, fmt.Errorf("column: %w", domain.ErrNotFound)
	}
	if err != nil {
		return domain.Column{}, fmt.Errorf("getOne error: %w", err)
	}

	return res, nil
}

func (c *columnRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Column, error) {
	query := `SELECT id,position,name,status,project_id FROM columns WHERE id = $1`

	return c.getOne(ctx, query, id)
}

func (c *columnRepository) Update(ctx context.Context, cls ...domain.Column) error {
	txn, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}
	query := `UPDATE columns SET position=$2,name=$3,status=$4,project_id=$5 WHERE id = $1`
	stmt, err := txn.Prepare(query)
	if err != nil {
		return fmt.Errorf("tx prepare: %w", err)
	}
	defer stmt.Close()

	for _, c := range cls {
		_, err = stmt.Exec(c.ID, c.Position, c.Name, c.Status, c.ProjectID)
		if err != nil {
			return fmt.Errorf("smt exec: %w", err)
		}
	}
	if err = stmt.Close(); err != nil {
		return fmt.Errorf("stm close: %w", err)
	}
	if err = txn.Commit(); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

func (c *columnRepository) Store(ctx context.Context, a *domain.Column) error {
	query := `INSERT INTO columns (position,name,status,project_id) VALUES ( $1, $2, $3, $4) RETURNING id`
	row := c.db.QueryRowContext(ctx, query, a.Position, a.Name, a.Status, a.ProjectID)
	err := row.Scan(&a.ID)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	return nil
}

func (c *columnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM columns WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}
