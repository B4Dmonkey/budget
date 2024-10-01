package app

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type MockQueries struct {
	db interface{}
}

func NewMockQueries(db any) *MockQueries {
	return &MockQueries{db: db}
}

func (m *MockQueries) GetDatum() string {
	return "Mock Datum"
}
func TestDatabase(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.Nil(t, err, "Error connecting to database")
	query := NewMockQueries(db)

	app := New(
		WithDbQueries(query),
	)

	app.Get("/", func(ctx Context) error {
		value := ctx.DB.(*MockQueries).GetDatum()
		return ctx.Send(value)
	})
}
