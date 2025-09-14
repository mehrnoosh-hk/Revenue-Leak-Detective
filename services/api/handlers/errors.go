package handlers

import "errors"

var (
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrHealthCheckFailed   = errors.New("health check failed")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotAlive            = errors.New("not alive")
	ErrNotReady            = errors.New("not ready")
)
