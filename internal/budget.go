package internal

import (
	"context"
	"io"
	"mime/multipart"
	"my-budget/database/orm"
	"time"
)

type Budget struct {
	ctx   context.Context
	query orm.Querier
}

func NewBudget(ctx context.Context, query orm.Querier) (*Budget, error) {
	// verify := Verifier{}
	// verify.That(query != nil, "Querier not provided")
	// verify.That(ctx != nil, "Context is nil")

	// if err := verify.Flush(); err != nil {
	// 	return nil, err
	// }

	return &Budget{ctx: ctx, query: query}, nil
}

func (b *Budget) GetTransactions(startDate time.Time, endDate time.Time) ([]orm.GetTransactionsInDateRangeRow, error) {
	// ? Validate the inputs
	// verify := Verifier{}
	// verify.That(startDate != nil, "Querier not provided")
	// verify.That(endDate != nil, "Context is nil")

	// if err := verify.Flush(); err != nil {
	// 	return nil, err
	// }

	return b.query.GetTransactionsInDateRange(
		b.ctx,
		orm.GetTransactionsInDateRangeParams{StartDate: startDate, EndDate: endDate},
	)
}

func (b *Budget) AddNewTransactionsFromDocument(header *multipart.FileHeader, file io.Reader) error {
	return nil
}