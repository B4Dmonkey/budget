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

func (m MockData) Render() (string, error) {
	template := m.Template()
	parsedTemplate, _ := mustache.ParseString(template)
	return parsedTemplate.Render(template, m)
}

func (m MockData) Template() string {
	return "{{StringField}}{{NumberField}}"
}

type MockDataSlice []MockData

func (m MockDataSlice) Template() string {
	return "{{#MockDataSlice}}{{StringField}}{{NumberField}}{{/MockDataSlice}}"
}

func (m MockDataSlice) Render() (string, error) {
	template := m.Template()
	parsedTemplate, _ := mustache.ParseString(template)
	return parsedTemplate.Render(template, map[string]MockDataSlice{"MockDataSlice": m})
}

// func TestRender(t *testing.T) {
// 	expected := "Foo64"
// 	mock_data := MockData{"Foo", 64}
// 	result, err := mock_data.Render()
// 	assert.Nil(t, err)
// 	assert.Equal(t, expected, result)
// }

func TestRendersMultipleDataPoints(t *testing.T) {
	expected := "Foo64Boo78"
	mock_slice_data := []MockData{
		{"Foo", 64},
		{"Boo", 78},
	}
	mock_data := MockDataSlice(mock_slice_data)
	result, err := mock_data.Render()
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRender(t *testing.T) {
	app := New()
	app.AddHandler(MethodGet, "/foo-boo", func(ctx Context) error {
		return ctx.Render(http.StatusOK)
	})

	tests := []struct {
		description string
		code        int
		method      string
		uri         string
		expected   string
	}{
		{
			description: "GET /foo-boo should return status 200 and render Foo64Boo78",
			code:        200,
			method:      "GET",
			uri:         "/foo-boo",
			expected:   "Foo64Boo78",
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
