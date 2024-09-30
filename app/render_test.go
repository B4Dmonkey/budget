package app

import (
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

func TestRender(t *testing.T) {
	expected := "Foo64"
	mock_data := MockData{"Foo", 64}
	result, err := mock_data.Render()
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

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
