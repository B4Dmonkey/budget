package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"my-budget/database/orm"
)

func TestConvertCurrencyIntToString(t *testing.T) {
	tests := map[string]struct {
		input    int64
		expected string
	}{
		"It converts 0 to 0.00":     {input: 0, expected: "0.00"},
		"It converts 100 to 1.23":   {input: 123, expected: "1.23"},
		"It converts -100 to -1.23": {input: -123, expected: "-1.23"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, ConvertCurrencyIntToString(tc.input))
		})
	}
}

func TestConvertCurrencyStringToIntOrNil(t *testing.T) {
	tests := map[string]struct {
		input     string
		expected  sql.NullInt64
		expectErr bool
	}{
		"It converts 0.00 to 0":              {input: "0.00", expected: sql.NullInt64{Int64: 0, Valid: true}},
		"It converts 1.23 to 123":            {input: "1.23", expected: sql.NullInt64{Int64: 123, Valid: true}},
		"It converts -1.23 to -123":          {input: "-1.23", expected: sql.NullInt64{Int64: -123, Valid: true}},
		"It returns empty string as invalid": {input: "", expected: sql.NullInt64{Valid: false}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := ConvertCurrencyStringToNullInt64(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func generateTestCSVData() *bytes.Buffer {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #", ""})
	writer.Write([]string{"", "01/01/2023", "Deposit", "100.50", "Credit", "100.50", ""})
	// writer.Write([]string{"", "01/02/2023", "Withdrawal", "-50.25", "Debit", "50.25", ""})
	// writer.Write([]string{"", "01/03/2023", "Deposit", "200.00", "Credit", "250.25", ""})
	// writer.Write([]string{"", "01/04/2023", "Withdrawal", "-75.00", "Debit", "175.25", ""})
	// writer.Write([]string{"", "01/05/2023", "Deposit", "150.75", "Credit", "326.00", ""})
	writer.Flush()
	return &buf
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateDocumentMeta(ctx context.Context, arg orm.CreateDocumentMetaParams) (orm.DocumentsMetum, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(orm.DocumentsMetum), args.Error(1)
}

func (m *MockDB) CreateTransaction(ctx context.Context, arg orm.CreateTransactionParams) error {
	return nil
	// args := m.Called(ctx, arg)
	// return args.Error(0)
}

func (m *MockDB) FindOneDocumentMeta(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockDB) GetPendingTransactions(ctx context.Context, documentID interface{}) ([]orm.Transaction, error) {
	args := m.Called(ctx, documentID)
	return args.Get(0).([]orm.Transaction), args.Error(1)
}

func (m *MockDB) GetTransactionsInDateRange(ctx context.Context, arg orm.GetTransactionsInDateRangeParams) ([]orm.Transaction, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]orm.Transaction), args.Error(1)
}

func TestProcessNewTransactions(t *testing.T) {
	mock_db := new(MockDB)
	mock_ctx := context.Background()
	mock_document_name := "test.csv"
	mock_document_id := "new-document-id"
	header := &multipart.FileHeader{Filename: mock_document_name}

	mock_csv_data := generateTestCSVData()
	mock_db.
		On("FindOneDocumentMeta", mock_ctx, header.Filename).
		Return("", sql.ErrNoRows)
	mock_db.
		On("CreateDocumentMeta", mock_ctx, orm.CreateDocumentMetaParams{
			Name:         header.Filename,
			PersistedLoc: header.Filename,
		}).
		Return(orm.DocumentsMetum{ID: mock_document_id}, nil)
	// mock_db.("CreateTransaction", mock_ctx, orm.CreateTransactionParams{})
	// mock.ExpectQuery(orm.FindOneDocumentMeta).
	// 	WithArgs(header.Filename).
	// 	WillReturnError(sql.ErrNoRows)

	// mock_document_row := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "persisted_loc"}).
	// 	AddRow("new-document-id", time.Now(), time.Now(), header.Filename, mock_document_name)

	// mock.ExpectQuery(orm.CreateDocumentMeta).
	// 	WithArgs(header.Filename, mock_document_name).
	// 	WillReturnRows(mock_document_row)

	err := ProcessNewTransactions(mock_db, mock_ctx, header, mock_csv_data)
	assert.Nil(t, err, "Expected an error")
}