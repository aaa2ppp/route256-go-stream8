package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"route256/cart/internal/model"
	"route256/cart/pkg/http/middleware"
)

type helper struct {
	w  http.ResponseWriter
	r  *http.Request
	lg *slog.Logger
}

func newHelper(w http.ResponseWriter, r *http.Request, op string) helper {
	return helper{
		w: w,
		r: r,
		lg: middleware.GetLoggerFromContextOrDefault(r.Context()).
			With("op", "handler#"+op),
	}
}

func (x helper) ctx() context.Context {
	return x.r.Context()
}

func (x helper) log() *slog.Logger {
	return x.lg
}

func (x helper) checkPOSTMethod() bool {
	if x.r.Method != http.MethodPost {
		x.log().Error("logical error: using for a non-POST request", "method", x.r.Method, "url", x.r.URL)
		x.writeError(errInternalError)
		return false
	}
	return true
}

func (x helper) getAuthToken() string {
	return x.r.Header.Get("X-Authtoken")
}

func (x helper) ReadBody() ([]byte, bool) {
	body, err := io.ReadAll(x.r.Body)
	if err != nil {
		x.log().Error("can't read request body", "error", err)
		x.writeError(errInternalError)
		return nil, false
	}
	return body, true
}

func (x helper) decodeBody(req any) bool {
	body, ok := x.ReadBody()
	if !ok {
		return false
	}
	if err := json.Unmarshal(body, req); err != nil {
		x.log().Debug("can't unmarshal request body", "error", err)
		x.writeError(errBadRequest)
		return false
	}
	return true
}

type Validator interface {
	Validate() error
}

func (x helper) validateRequest(req any) bool {
	if req2, ok := req.(Validator); !ok {
		x.log().Error("request not implement Validator", "reqType", reflect.TypeOf(req))
		x.writeError(errInternalError)
		return false
	} else if err := req2.Validate(); err != nil {
		x.log().Debug("can't validate request", "error", err)
		x.writeError(errBadRequest)
		return false
	}
	return true
}

func (x helper) decodeBodyAndValidateRequest(req any) bool {
	return x.decodeBody(req) && x.validateRequest(req)
}

func (x helper) writeResponse(status int, resp any) {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		x.log().Error("can't marshal response", "error", err)
		x.writeError(errInternalError)
		return
	}

	x.w.Header().Add("content-type", "application/json")
	x.w.WriteHeader(status)

	if _, err := x.w.Write(jsonData); err != nil {
		x.log().Error("can't write response", "error", err)
		return
	}
}

func (x helper) writeError(err error) {
	var httpErr *httpError
	switch {
	case errors.As(err, &httpErr):
	case errors.Is(err, model.ErrUnauthorized):
		httpErr = errUnauthorized
	case errors.Is(err, model.ErrNotFound):
		httpErr = errNotFound
	case errors.Is(err, model.ErrPreconditionFailed):
		httpErr = errPreconditionFailed
	case errors.Is(err, model.ErrInternalError):
		httpErr = errInternalError
	default:
		x.log().Warn("unhandled error", "error", err)
		httpErr = errInternalError
	}
	http.Error(x.w, httpErr.Message, httpErr.Code)
}
