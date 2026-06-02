package service

import "errors"

var (
	notFoundErr = errors.New("Not found")
	internalErr = errors.New("Internal server error")
)
