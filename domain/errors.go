package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen.
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists.
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists.
	ErrConflict = errors.New("item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid.
	ErrBadParamInput = errors.New("given param is not valid")
)
