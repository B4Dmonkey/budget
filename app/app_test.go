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

	app.AddHandler(http.MethodGet, "/foo", func(ctx Context) error {
		return ctx.Status(218)
	})
	app.Get("/", func(ctx Context) error {
		return ctx.Status(http.StatusOK)
	})

	app.Post("/", func(ctx Context) error {
		return ctx.Status(http.StatusCreated)
	})

	tests := []struct {
		description string
		code        int
		method      string
		uri         string
	}{
		{
			description: "GET /foo should return 218",
			code:        218,
			method:      "GET",
			uri:         "/foo",
		},
		{
			description: "GET / should return 200",
			code:        http.StatusOK,
			method:      "GET",
			uri:         "/",
		},
		{
			description: "POST / should return 201",
			code:        http.StatusCreated,
			method:      "POST",
			uri:         "/",
		},
	}

	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.uri, nil)
		w := httptest.NewRecorder()
		app.mux.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code, test.description)
	}
}

func TestAppListen(t *testing.T) {
	app := New()
	go func() {
		err := app.Listen(":8080")
		assert.NotNil(t, err, "Listen should return an error")
	}()
}
