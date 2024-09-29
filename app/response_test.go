package app

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_Status(t *testing.T) {
	rr := httptest.NewRecorder()

	response := &Response{w: rr}

	statusCode := http.StatusOK
	err := response.Status(statusCode)
	assert.Nil(t, err)
	assert.Equal(t, statusCode, rr.Code)
}
