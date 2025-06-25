package storeerrors

import "errors"

var (
	ErrTaskCompleted = errors.New("task already complete")
	ErrTaskNotFound  = errors.New("task not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user aldready exists")
)
