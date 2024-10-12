package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMiddleware(t *testing.T) {
	sampleHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	sampleMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Middleware"))
			next.ServeHTTP(w, r)
		}
	}
	wrappedHandler := WithMiddleware(sampleHandler, sampleMiddleware)
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.Nil(t, err, "Error creating request")
	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Status code not OK")
	assert.Equal(t, "Middleware", rr.Body.String(), "Response body not as expected")
}
