package service

import "errors"

var (
	NotFoundErr         = errors.New("Not found")
	internalErr         = errors.New("Internal server error")
	incorrectTestSubmit = errors.New("Incorrect test submit")
)
