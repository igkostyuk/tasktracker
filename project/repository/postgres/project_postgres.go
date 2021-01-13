package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
)

type projectRepository struct {
	db *sql.DB
}

// New will create new a projectRepository object representation of domain.ProjectRepository interface.
func New(db *sql.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (p *projectRepository) Fetch(ctx context.Context) ([]domain.Project, error) {
	query := `SELECT id,name,description FROM projects`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	result := make([]domain.Project, 0)
	for rows.Next() {
		t := domain.Project{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
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

func (p *projectRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Project, error) {
	row := p.db.QueryRowContext(ctx, query, args...)
	res := domain.Project{}
	err := row.Scan(&res.ID, &res.Name, &res.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Project{}, fmt.Errorf("project: %w", domain.ErrNotFound)
	}
	if err != nil {
		return domain.Project{}, fmt.Errorf("getOne error: %w", err)
	}

	return res, nil
}

func (p *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Project, error) {
	query := `SELECT id,name,description FROM projects WHERE id = $1`

	return p.getOne(ctx, query, id)
}

func (p *projectRepository) Update(ctx context.Context, pr *domain.Project) error {
	query := `UPDATE projects SET name = $2,description = $3 FROM projects WHERE id = $1`
	_, err := p.db.ExecContext(ctx, query, pr.ID, pr.Name, pr.Description)
	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	return nil
}

func (p *projectRepository) Store(ctx context.Context, a *domain.Project) error {
	query := `INSERT INTO projects ( name, description) VALUES ($1, $2) RETURNING id`
	row := p.db.QueryRowContext(ctx, query, a.Name, a.Description)
	err := row.Scan(&a.ID)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	return nil
}

func (p *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	return nil
}
