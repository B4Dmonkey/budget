package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"

	"my-budget/database/orm"
)

func generateTestCSVData() *bytes.Buffer {
	// ToDo: This should be randomly generated
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	headers := []string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #", ""}
	row := []string{"DEBIT", "01/01/2023", "MTA*NYCT PAYGO NEW YORK NY                   09/23", "100.50", "DEBIT", "100.50", ""}
	if err := writer.Write(headers); err != nil {
		panic(err)
	}
	if err := writer.Write(row); err != nil {
		panic(err)
	}
	writer.Flush() // Ensure all data is written to the buffer
	if err := writer.Error(); err != nil {
		panic(err)
	}
	return &buf
}

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

func TestGetTransactions(t *testing.T) {
	tests := map[string]struct {
		db        orm.Querier
		ctx       context.Context
		expectErr bool
	}{
		"It fails when db is nil":  {db: nil, ctx: context.Background(), expectErr: true},
		"It fails when ctx is nil": {db: new(MockDB), ctx: nil, expectErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := GetTransactions(tc.db, tc.ctx)
			if tc.expectErr {
				assert.Error(t, err)
			}
		})
	}
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
	assert.NoError(t, err, "Expected nil error")
}
