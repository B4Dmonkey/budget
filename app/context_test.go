package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestContext(t *testing.T) {
	app := New()

	app.Get("/", func(ctx Context) error {
		assert.NotNil(t, ctx.Req, "ctx.Req should not be nil")
		assert.NotNil(t, ctx.Req.Context(), "ctx.Res should not be nil")
		return nil
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.mux.ServeHTTP(w, req)
}