package router_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/domain"
	mocks "github.com/igkostyuk/tasktracker/domain/mock"
	"github.com/igkostyuk/tasktracker/internal/web"
	projectDelivery "github.com/igkostyuk/tasktracker/project/delivery/http"
	helper "github.com/matryer/is"
)

var validUUIDString = "177ef0d8-6630-11ea-b69a-0242ac130003"

func TestFetch(t *testing.T) {
	t.Run("returns Projects", func(t *testing.T) {
		is := helper.New(t)
		want := []domain.Project{{Name: "1", Description: "testDescription"}, {Name: "2", Description: "testDescription2"}}
		// nolint:exhaustivestruct
		mockedProjectUsecase := &mocks.ProjectUsecaseMock{
			FetchFunc: func(ctx context.Context) ([]domain.Project, error) {
				return want, nil
			},
		}
		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

		var got []domain.Project
		err := json.NewDecoder(response.Body).Decode(&got)
		is.NoErr(err)

		is.Equal(response.Code, http.StatusOK)
		is.Equal(want, got)
		is.Equal(len(mockedProjectUsecase.FetchCalls()), 1)
	})
	t.Run("returns error", func(t *testing.T) {
		is := helper.New(t)
		mockError := domain.ErrInternalServerError

		// nolint:exhaustivestruct
		mockedProjectUsecase := &mocks.ProjectUsecaseMock{
			FetchFunc: func(ctx context.Context) ([]domain.Project, error) {
				return nil, mockError
			},
		}
		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		is.NoErr(err)
		checkError(t, mockedProjectUsecase, request, http.StatusInternalServerError, mockError.Error())
		is.Equal(len(mockedProjectUsecase.FetchCalls()), 1)
	})
}

func TestFetchColumns(t *testing.T) {
	is := helper.New(t)
	want := []domain.Column{{Name: "1", Position: 1, Status: "testStatus"}}
	var calledUUID uuid.UUID
	// nolint:exhaustivestruct
	mockedProjectUsecase := &mocks.ProjectUsecaseMock{
		FetchColumnsFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
			calledUUID = id

			return want, nil
		},
	}
	path := fmt.Sprintf("/%s/columns", validUUIDString)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

	var got []domain.Column
	err := json.NewDecoder(response.Body).Decode(&got)
	is.NoErr(err)

	is.Equal(response.Code, http.StatusOK)
	is.Equal(want, got)
	is.Equal(calledUUID.String(), validUUIDString)
	is.Equal(len(mockedProjectUsecase.FetchColumnsCalls()), 1)
}

func TestFetchColumnsErrors(t *testing.T) {
	tt := []struct {
		name      string
		path      string
		code      int
		message   string
		mockError error
	}{
		{
			name:      "500",
			path:      fmt.Sprintf("/%s/columns", validUUIDString),
			code:      http.StatusInternalServerError,
			message:   domain.ErrInternalServerError.Error(),
			mockError: domain.ErrInternalServerError,
		},
		{
			name:      "404 invalid path",
			path:      fmt.Sprintf("/%s/columns", "invalidUUID"),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: nil,
		},
		{
			name:      "404 not found project",
			path:      fmt.Sprintf("/%s/columns", validUUIDString),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: domain.ErrNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			var calledUUID uuid.UUID
			// nolint:exhaustivestruct
			mockedProjectUsecase := &mocks.ProjectUsecaseMock{
				FetchColumnsFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Column, error) {
					calledUUID = id

					return nil, tc.mockError
				},
			}
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tc.path, nil)
			is.NoErr(err)
			checkError(t, mockedProjectUsecase, request, tc.code, tc.message)

			if tc.mockError != nil {
				is.Equal(len(mockedProjectUsecase.FetchColumnsCalls()), 1)
				is.Equal(calledUUID.String(), validUUIDString)

				return
			}
			is.Equal(len(mockedProjectUsecase.FetchColumnsCalls()), 0)
		})
	}
}

func TestFetchTasks(t *testing.T) {
	is := helper.New(t)
	want := []domain.Task{{Name: "1", Position: 1, Description: "testDescription"}}

	var calledUUID uuid.UUID
	// nolint:exhaustivestruct
	mockedProjectUsecase := &mocks.ProjectUsecaseMock{
		FetchTasksFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
			calledUUID = id

			return want, nil
		},
	}
	path := fmt.Sprintf("/%s/tasks", validUUIDString)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

	var got []domain.Task
	err := json.NewDecoder(response.Body).Decode(&got)
	is.NoErr(err)

	is.Equal(response.Code, http.StatusOK)
	is.Equal(want, got)
	is.Equal(calledUUID.String(), validUUIDString)
	is.Equal(len(mockedProjectUsecase.FetchTasksCalls()), 1)
}

func TestFetchTasksErrors(t *testing.T) {
	tt := []struct {
		name      string
		path      string
		code      int
		message   string
		mockError error
	}{
		{
			name:      "500",
			path:      fmt.Sprintf("/%s/tasks", validUUIDString),
			code:      http.StatusInternalServerError,
			message:   domain.ErrInternalServerError.Error(),
			mockError: domain.ErrInternalServerError,
		},
		{
			name:      "404 invalid path",
			path:      fmt.Sprintf("/%s/tasks", "invalidUUID"),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: nil,
		},
		{
			name:      "404 not found project",
			path:      fmt.Sprintf("/%s/tasks", validUUIDString),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: domain.ErrNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			var calledUUID uuid.UUID
			// nolint:exhaustivestruct
			mockedProjectUsecase := &mocks.ProjectUsecaseMock{
				FetchTasksFunc: func(ctx context.Context, id uuid.UUID) ([]domain.Task, error) {
					calledUUID = id

					return nil, tc.mockError
				},
			}
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tc.path, nil)
			is.NoErr(err)
			checkError(t, mockedProjectUsecase, request, tc.code, tc.message)

			if tc.mockError != nil {
				is.Equal(len(mockedProjectUsecase.FetchTasksCalls()), 1)
				is.Equal(calledUUID.String(), validUUIDString)

				return
			}
			is.Equal(len(mockedProjectUsecase.FetchTasksCalls()), 0)
		})
	}
}

func TestGetByID(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	want := domain.Project{Name: "testName", Description: "testDescription"}
	var calledUUID uuid.UUID
	// nolint:exhaustivestruct
	mockedProjectUsecase := &mocks.ProjectUsecaseMock{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
			calledUUID = id

			return want, nil
		},
	}
	path := fmt.Sprintf("/%s", validUUIDString)
	request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

	var got domain.Project
	err := json.NewDecoder(response.Body).Decode(&got)
	is.NoErr(err)

	is.Equal(response.Code, http.StatusOK)
	is.Equal(want, got)
	is.Equal(calledUUID.String(), validUUIDString)
	is.Equal(len(mockedProjectUsecase.GetByIDCalls()), 1)
}

func TestGetByIDErrors(t *testing.T) {
	tt := []struct {
		name      string
		path      string
		code      int
		message   string
		mockError error
	}{
		{
			name:      "500",
			path:      fmt.Sprintf("/%s", validUUIDString),
			code:      http.StatusInternalServerError,
			message:   domain.ErrInternalServerError.Error(),
			mockError: domain.ErrInternalServerError,
		},
		{
			name:      "404 invalid path",
			path:      fmt.Sprintf("/%s", "invalidUUID"),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: nil,
		},
		{
			name:      "404 not found project",
			path:      fmt.Sprintf("/%s", validUUIDString),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: domain.ErrNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			var calledUUID uuid.UUID
			// nolint:exhaustivestruct
			mockedProjectUsecase := &mocks.ProjectUsecaseMock{
				GetByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Project, error) {
					calledUUID = id

					return domain.Project{}, tc.mockError
				},
			}
			request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tc.path, nil)
			is.NoErr(err)
			checkError(t, mockedProjectUsecase, request, tc.code, tc.message)

			if tc.mockError != nil {
				is.Equal(len(mockedProjectUsecase.GetByIDCalls()), 1)
				is.Equal(calledUUID.String(), validUUIDString)

				return
			}
			is.Equal(len(mockedProjectUsecase.GetByIDCalls()), 0)
		})
	}
}

func TestStore(t *testing.T) {
	is := helper.New(t)
	// nolint:exhaustivestruct
	want := domain.Project{Name: "testName", Description: "testDescription"}

	// nolint:exhaustivestruct
	mockedProjectUsecase := &mocks.ProjectUsecaseMock{
		StoreFunc: func(ctx context.Context, pr *domain.Project) error {
			return nil
		},
	}
	jsonData, err := json.Marshal(want)
	is.NoErr(err)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewBuffer(jsonData))
	is.NoErr(err)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

	is.Equal(response.Code, http.StatusOK)
	var got domain.Project
	err = json.NewDecoder(response.Body).Decode(&got)
	is.NoErr(err)

	is.Equal(want, got)
	is.Equal(len(mockedProjectUsecase.StoreCalls()), 1)
}

//nolint:funlen
func TestStoreErrors(t *testing.T) {
	tt := []struct {
		name      string
		project   string
		code      int
		message   string
		mockError error
	}{
		{
			name:    "400",
			project: `{"name":"testName"}`,
			code:    http.StatusBadRequest,
			message: "validation: Key: 'Project.Description'" +
				" Error:Field validation for 'Description' failed on the 'required' tag",
			mockError: nil,
		},
		{
			name:      "400",
			project:   `{"description":"testDescription"}`,
			code:      http.StatusBadRequest,
			message:   "validation: Key: 'Project.Name' Error:Field validation for 'Name' failed on the 'required' tag",
			mockError: nil,
		},
		{
			name:      "409",
			project:   `{"name":"testName","description":"testDescription"}`,
			code:      http.StatusConflict,
			message:   domain.ErrConflict.Error(),
			mockError: domain.ErrConflict,
		},
		{
			name:      "422",
			project:   ``,
			code:      http.StatusUnprocessableEntity,
			message:   "EOF",
			mockError: nil,
		},
		{
			name:      "500",
			project:   `{"name":"testName","description":"testDescription"}`,
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: domain.ErrNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			// nolint:exhaustivestruct
			mockedProjectUsecase := &mocks.ProjectUsecaseMock{
				StoreFunc: func(ctx context.Context, pr *domain.Project) error {
					return tc.mockError
				},
			}

			request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader(tc.project))
			is.NoErr(err)
			request.Header.Set("Content-Type", "application/json")
			checkError(t, mockedProjectUsecase, request, tc.code, tc.message)

			if tc.mockError != nil {
				is.Equal(len(mockedProjectUsecase.StoreCalls()), 1)

				return
			}
			is.Equal(len(mockedProjectUsecase.StoreCalls()), 0)
		})
	}
}

func TestDelete(t *testing.T) {
	is := helper.New(t)
	var calledUUID uuid.UUID
	// nolint:exhaustivestruct
	mockedProjectUsecase := &mocks.ProjectUsecaseMock{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			calledUUID = id

			return nil
		},
	}
	path := fmt.Sprintf("/%s", validUUIDString)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, path, nil)
	is.NoErr(err)
	response := httptest.NewRecorder()

	projectDelivery.New(mockedProjectUsecase).ServeHTTP(response, request)

	is.Equal(response.Code, http.StatusNoContent)
	is.Equal(calledUUID.String(), validUUIDString)
	is.Equal(len(mockedProjectUsecase.DeleteCalls()), 1)
}

func TestDeleteErrors(t *testing.T) {
	tt := []struct {
		name      string
		path      string
		code      int
		message   string
		mockError error
	}{
		{
			name:      "500",
			path:      fmt.Sprintf("/%s", validUUIDString),
			code:      http.StatusInternalServerError,
			message:   domain.ErrInternalServerError.Error(),
			mockError: domain.ErrInternalServerError,
		},
		{
			name:      "404 invalid path",
			path:      fmt.Sprintf("/%s", "invalidUUID"),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: nil,
		},
		{
			name:      "404 not found project",
			path:      fmt.Sprintf("/%s", validUUIDString),
			code:      http.StatusNotFound,
			message:   domain.ErrNotFound.Error(),
			mockError: domain.ErrNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := helper.New(t)
			var calledUUID uuid.UUID
			// nolint:exhaustivestruct
			mockedProjectUsecase := &mocks.ProjectUsecaseMock{
				DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
					calledUUID = id

					return tc.mockError
				},
			}
			request, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, tc.path, nil)
			is.NoErr(err)
			checkError(t, mockedProjectUsecase, request, tc.code, tc.message)

			if tc.mockError != nil {
				is.Equal(len(mockedProjectUsecase.DeleteCalls()), 1)
				is.Equal(calledUUID.String(), validUUIDString)

				return
			}
			is.Equal(len(mockedProjectUsecase.DeleteCalls()), 0)
		})
	}
}

func checkError(t *testing.T, mockedUsecase domain.ProjectUsecase, request *http.Request, code int, message string) {
	t.Helper()
	is := helper.New(t)

	response := httptest.NewRecorder()

	projectDelivery.New(mockedUsecase).ServeHTTP(response, request)

	want := web.HTTPError{Code: code, Message: message}

	var got web.HTTPError
	err := json.NewDecoder(response.Body).Decode(&got)

	is.NoErr(err)
	is.Equal(response.Code, code)
	is.Equal(want, got)
}
