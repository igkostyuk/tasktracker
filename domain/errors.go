package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen.
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists.
	ErrNotFound = errors.New("item is not found")
	// ErrConflict will throw if the current action already exists.
	ErrConflict = errors.New("item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid.
	ErrBadParamInput = errors.New("given param is not valid")
	// ErrLastColumn will throw if trying to delete the last column.
	ErrLastColumn = errors.New("the last column cannot be deleted")
	// ErrUnique will throw if column name or status not unique for project.
	ErrUnique = errors.New("must be unique")
)
