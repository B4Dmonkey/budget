package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cbroglie/mustache"
	"github.com/stretchr/testify/assert"
)

type MockData struct {
	StringField string
	NumberField float64
}

func (m MockData) Binding() interface{} {
	return nil
}

func (m MockData) Template() (*mustache.Template, error) {
	return nil, nil
}

func (m MockData) Render() (string, error) {
	return "Don DADA!!!", nil
}

type MockDataSlice []MockData

func (m MockDataSlice) Binding() interface{} {
	return map[string]interface{}{
		"MockDataSlice": m,
	}
}

func (m MockDataSlice) Template() (*mustache.Template, error) {
	template := "{{#MockDataSlice}}{{StringField}}{{NumberField}}{{/MockDataSlice}}"
	return mustache.ParseString(template)
}

func TestRender(t *testing.T) {
	app := New()
	app.Get("/foo-boo", func(ctx Context) error {
		mock_slice_data := []MockData{
			{"Foo", 64},
			{"Boo", 78},
		}
		mock_data := MockDataSlice(mock_slice_data)
		return ctx.Render(http.StatusOK, mock_data)
	})

	app.Get("/overwritten-render", func(ctx Context) error {
		mock_data := new(MockData)
		return ctx.Render(http.StatusOK, mock_data)
	})

	tests := []struct {
		description string
		code        int
		method      string
		uri         string
		expected    string
	}{
		{
			description: "GET /foo-boo should return status 200 and render Foo64Boo78",
			code:        200,
			method:      "GET",
			uri:         "/foo-boo",
			expected:    "Foo64Boo78",
		},
		{
			description: "GET /overwritten-render should return status 200 and render Don DADA!!!",
			code:        200,
			method:      "GET",
			uri:         "/overwritten-render",
			expected:    "Don DADA!!!",
		},
	}
	for _, test := range tests {
		req, _ := http.NewRequest(test.method, test.uri, nil)
		w := httptest.NewRecorder()
		app.mux.ServeHTTP(w, req)
		assert.Equal(t, test.code, w.Code, test.description)
		assert.Contains(t, w.Body.String(), test.expected, test.description)
	}
}
