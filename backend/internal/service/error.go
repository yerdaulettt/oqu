package service

import "errors"

var (
	NotFoundErr         = errors.New("Not found")
	internalErr         = errors.New("Internal server error")
	incorrectTestSubmit = errors.New("Incorrect test submit")
	IncorrectPassword   = errors.New("Incorrect password")
	AlreadyCompleted    = errors.New("Already completed")
	UsernameErr         = errors.New("Username exists")
	IncorrectRole       = errors.New("Incorrect role")
)
