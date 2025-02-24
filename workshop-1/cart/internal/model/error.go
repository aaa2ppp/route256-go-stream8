package model

import "errors"

var (
	ErrInternalError      = errors.New("internal error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrNotFound           = errors.New("not found")
	ErrPreconditionFailed = errors.New("precondition failed")
)
