package domain

import "errors"

const (
	NotAdded     = "It was not added."
	NotFound     = "It was not found."
	AlreadyExist = "It already exists."
)

var (
	ErrNotAdded     = errors.New(NotAdded)
	ErrNotFound     = errors.New(NotFound)
	ErrAlreadyExist = errors.New(AlreadyExist)
)
