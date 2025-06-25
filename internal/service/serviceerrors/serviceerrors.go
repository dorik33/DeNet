package serviceerrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyDone   = errors.New("task already completed")
	ErrInvalidPassword   = errors.New("invalid password")
)
