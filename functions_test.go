package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"my-budget/database/orm"
	"my-budget/internal/errors"
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

func setupTestGetTransactions(testQuery *Domain) {
	filename := "Chase Activity Sept 27.CSV"

	document, err := testQuery.q.CreateDocumentMeta(testQuery.ctx, orm.CreateDocumentMetaParams{
		Name:         filename,
		PersistedLoc: filename,
		PublishingDate: sql.NullTime{
			Time:  time.Date(2024, time.September, 15, 0, 0, 0, 0, time.Local),
			Valid: true,
		},
	})

	if err != nil {
		testQuery.t.Fatalf("Failed to create document meta: %s", err)
	}

	err = testQuery.q.CreateTransaction(testQuery.ctx, orm.CreateTransactionParams{
		DocumentID:  document.ID,
		Details:     "Test Transaction Details",
		PostingDate: time.Date(2024, time.September, 27, 0, 0, 0, 0, time.Local),
		Amount:      123,
		Type:        "Credit",
		Balance:     sql.NullInt64{Int64: 123, Valid: true},
	})

	if err != nil {
		testQuery.t.Fatalf("Failed to create transaction durning setup: %s", err)
	}

	if err := testQuery.tx.Commit(); err != nil {
		testQuery.t.Fatalf("Failed to commit transaction: %s", err)
	}
}
func TestGetTransactions(t *testing.T) {
	dt := NewDomainTest(t)
	defer dt.teardown()

	setupTestGetTransactions(dt)

	mock_ctx := context.Background()
	start_date := time.Date(2024, time.September, 1, 0, 0, 0, 0, time.Local)
	end_date := time.Date(2024, time.September, 30, 23, 59, 59, 999999999, time.Local)

	tests := map[string]struct {
		db        orm.Querier
		ctx       context.Context
		expectErr bool
	}{
		"It fails when db is nil":            {db: nil, ctx: mock_ctx, expectErr: true},
		"It fails when ctx is nil":           {db: dt.q, ctx: nil, expectErr: true},
		"It returns a slice of Transactions": {db: dt.q, ctx: mock_ctx, expectErr: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			results, err := GetTransactions(tc.db, tc.ctx, start_date, end_date)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, results, "Expected 1 result") // ? technically not nil but yea empty list
			}
		})
	}
}

func TestProcessNewTransactions(t *testing.T) {
	dt := NewDomainTest(t)
	defer dt.teardown()
	setupTestGetTransactions(dt)

	mock_document_name := "Chase Activity Sept 27.CSV"
	header := &multipart.FileHeader{Filename: mock_document_name}

	mock_csv_data, err := os.Open("testdata/Chase Activity Sept 27.CSV")
	if err != nil {
		t.Fatalf("failed to open test CSV file: %v", err)
	}
	defer mock_csv_data.Close()

	tests := map[string]struct {
		db              orm.Querier
		ctx             context.Context
		header          *multipart.FileHeader
		file            io.Reader
		expectErr       bool
		expectedErrType error
	}{
		"It fails when db is nil": {
			db: nil, ctx: dt.ctx, header: header, file: mock_csv_data, expectErr: true, expectedErrType: &errors.VerificationError{},
		},
		"It fails when ctx is nil": {
			db: dt.q, ctx: nil, header: header, file: mock_csv_data, expectErr: true, expectedErrType: &errors.VerificationError{},
		},
		"It fails when header is nil": {
			db: dt.q, ctx: dt.ctx, header: nil, file: mock_csv_data, expectErr: true, expectedErrType: &errors.VerificationError{},
		},
		"It fails when header filename is incorrect": {
			db:              dt.q,
			ctx:             dt.ctx,
			header:          &multipart.FileHeader{Filename: "test.csv"},
			file:            mock_csv_data,
			expectErr:       true,
			expectedErrType: &errors.VerificationError{},
		},
		"It fails when file is nil": {
			db: dt.q, ctx: dt.ctx, header: header, file: nil, expectErr: true, expectedErrType: &errors.VerificationError{},
		},
		"It fails when extracting date part": {
			db:              dt.q,
			ctx:             dt.ctx,
			header:          &multipart.FileHeader{Filename: "chase activity test.csv"},
			file:            mock_csv_data,
			expectErr:       true,
			expectedErrType: &errors.ParseError{},
		},
		"It does the thing": {db: dt.q, ctx: dt.ctx, header: header, file: mock_csv_data},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := ProcessNewTransactions(tc.db, tc.ctx, tc.header, tc.file)
			if tc.expectErr {
				assert.Error(t, err)
				assert.IsType(t, tc.expectedErrType, err, fmt.Sprintf("Received Error: %v", err))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
