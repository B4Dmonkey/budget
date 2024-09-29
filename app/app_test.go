package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app := New()

	assert.NotNil(t, app, "app should not be nil")
	assert.NotNil(t, app.mux, "app.mux should not be nil")
}

func TestAppHTTPMethods(t *testing.T) {
	app := New()

	app.Get("/", func(ctx Context) error {
		return ctx.Res.Status(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	app.mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "should return 200")
}
