package app

import (
	"context"
	"net/http"
)

type Request struct {
	r *http.Request
}

func (r *Request) Context() context.Context {
	return r.r.Context()
}