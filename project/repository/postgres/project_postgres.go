package postgres

import (
	"context"
	"database/sql"

	"github.com/igkostyuk/tasktracker/domain"
)

type projectRepository struct {
	conn *sql.DB
}

func New(conn *sql.DB) domain.ProjectRepository {
	return &projectRepository{conn: conn}
}

func (p *projectRepository) Fetch(ctx context.Context) ([]domain.Project, error) {
	query := `SELECT id,name,description FROM projects`
	rows, err := p.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result := make([]domain.Project, 0)
	for rows.Next() {
		t := domain.Project{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (p *projectRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Project, error) {
	row := p.conn.QueryRowContext(ctx, query, args...)
	res := domain.Project{}
	err := row.Scan(&res.ID, &res.Name, &res.Description)

	return res, err
}

func (p *projectRepository) GetByID(ctx context.Context, id string) (domain.Project, error) {
	query := `SELECT id,name,description FROM projects WHERE id = $1`

	return p.getOne(ctx, query, id)
}

func (p *projectRepository) Update(ctx context.Context, pr *domain.Project) error {
	query := `UPDATE projects SET name=$2,description=$3 FROM projects WHERE id = $1`
	_, err := p.conn.ExecContext(ctx, query, pr.ID, pr.Name, pr.Description)

	return err
}

func (p *projectRepository) Store(ctx context.Context, a *domain.Project) error {
	query := `INSERT INTO projects (name, description) VALUES ( $1, $2) RETURNING id`
	row := p.conn.QueryRowContext(ctx, query, a.Name, a.Description)
	err := row.Scan(&a.ID)
	return err
}

func (p *projectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := p.conn.ExecContext(ctx, query, id)

	return err
}
