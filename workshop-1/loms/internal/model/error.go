package model

import "errors"

var (
	ErrInternalError      = errors.New("internal error")
	ErrNotFound           = errors.New("not found")
	ErrPreconditionFailed = errors.New("precondition failed")
)
