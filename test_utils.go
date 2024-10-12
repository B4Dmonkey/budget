package main

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"my-budget/database/orm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateDocumentMeta(ctx context.Context, arg orm.CreateDocumentMetaParams) (orm.DocumentsMetum, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(orm.DocumentsMetum), args.Error(1)
}

func (m *MockDB) CreateTransaction(ctx context.Context, arg orm.CreateTransactionParams) error {
	return nil
}

func (m *MockDB) FindOneDocumentMeta(ctx context.Context, args orm.FindOneDocumentMetaParams) (string, error) {
	// args := m.Called(ctx, name)
	// return args.String(0), args.Error(1)
	return "", nil
}

func (m *MockDB) GetPendingTransactions(ctx context.Context, documentID interface{}) ([]orm.Transaction, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]orm.Transaction), args.Error(1)
}

func (m *MockDB) GetTransactionsInDateRange(
	ctx context.Context,
	arg orm.GetTransactionsInDateRangeParams,
) ([]orm.Transaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]orm.Transaction), args.Error(1)
}

func setupTestServer(t *testing.T) *httptest.Server {
	// ctx := context.Background()
	app := New()
	// assert.Equal(t, "test_address", app.server.Addr)
	ts := httptest.NewServer(app.mux)
	return ts
}
