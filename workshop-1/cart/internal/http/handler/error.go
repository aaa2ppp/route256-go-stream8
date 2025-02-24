package handler

import "fmt"

type httpError struct {
	Code    int
	Message string
}

func (e *httpError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

var (
	errInternalError      = &httpError{500, "internal error"}
	errBadRequest         = &httpError{400, "bad request"}
	errUnauthorized       = &httpError{401, "unauthorized"}
	errNotFound           = &httpError{404, "not found"}
	errPreconditionFailed = &httpError{412, "precondition failed"}
)
