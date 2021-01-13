package postgres_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
	projectRepository "github.com/igkostyuk/tasktracker/project/repository/postgres"
	helper "github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := helper.New(t)

	mockProjects := []domain.Project{
		{Name: "TestName1", Description: "testDescription1"},
		{Name: "TestName2", Description: "testDescription2"},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(mockProjects[0].ID, mockProjects[0].Name, mockProjects[0].Description).
		AddRow(mockProjects[1].ID, mockProjects[1].Name, mockProjects[1].Description)

	query := `SELECT id,name,description FROM projects`

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(query).WillReturnRows(rows)
		projects, err := projectRepository.New(db).Fetch(context.TODO())
		is.NoErr(err)
		is.Equal(mockProjects, projects)
	})
	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("some error"))
		_, err = projectRepository.New(db).Fetch(context.TODO())
		is.True(err != nil)
	})
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)
	mockProject := domain.Project{Name: "TestName1", Description: "testDescription1"}

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(mockProject.ID, mockProject.Name, mockProject.Description)

	query := `SELECT id,name,description FROM projects WHERE id = $1`
	var id uuid.UUID

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(id).WillReturnRows(rows)
		project, err := projectRepository.New(db).GetByID(context.TODO(), id)
		is.NoErr(err)
		is.Equal(mockProject, project)
	})
	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("some error"))
		_, err = projectRepository.New(db).GetByID(context.TODO(), id)
		is.True(err != nil)
	})

}
func TestUpdate(t *testing.T) {
	is := helper.New(t)
	pr := &domain.Project{Name: "TestName1", Description: "testDescription1"}

	query := `UPDATE projects SET name = $2,description = $3 FROM projects WHERE id = $1`

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(pr.ID, pr.Name, pr.Description).WillReturnResult(sqlmock.NewResult(1, 1))

		err = projectRepository.New(db).Update(context.TODO(), pr)
		is.NoErr(err)
	})
	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(fmt.Errorf("some error"))

		err = projectRepository.New(db).Update(context.TODO(), pr)
		is.True(err != nil)
	})
}

func TestStore(t *testing.T) {
	is := helper.New(t)

	pr := &domain.Project{Name: "TestName1", Description: "testDescription1"}
	query := `INSERT INTO projects ( name, description) VALUES ($1, $2) RETURNING id`
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pr.Name, pr.Description).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		err = projectRepository.New(db).Store(context.TODO(), pr)
		is.NoErr(err)
		is.Equal(pr.ID, id)
	})
	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(fmt.Errorf("some error"))

		err = projectRepository.New(db).Store(context.TODO(), pr)
		is.True(err != nil)
	})
}

func TestDelete(t *testing.T) {
	is := helper.New(t)

	id := uuid.New()
	query := `DELETE FROM projects WHERE id = $1`

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
		err = projectRepository.New(db).Delete(context.TODO(), id)
		is.NoErr(err)
	})
	t.Run("error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		is.NoErr(err)
		mock.ExpectExec(regexp.QuoteMeta(query)).WillReturnError(fmt.Errorf("some error"))
		err = projectRepository.New(db).Delete(context.TODO(), id)
		is.True(err != nil)
	})
}
