package internal

import (
	"context"
	"database/sql"
	stdErrs "errors"
	"fmt"
	"io"
	"log"
	"my-budget/database/orm"
	"my-budget/internal/errors"
	"strings"
	"time"
)

type Budget struct {
	ctx  context.Context
	conn *sql.DB
}

func NewBudget(ctx context.Context, conn *sql.DB) (*Budget, error) {
	verify := errors.Verifier{}
	verify.That(conn != nil, "Querier not provided")
	verify.That(ctx != nil, "Context is nil")

	if err := verify.Flush(); err != nil {
		return nil, err
	}

	return &Budget{ctx: ctx, conn: conn}, nil
}

// func (b *Budget) GetTransactions(startDate time.Time, endDate time.Time) ([]orm.GetTransactionsInDateRangeRow, error) {
// 	// ? Validate the inputs
// 	// verify := Verifier{}
// 	// verify.That(startDate != nil, "Querier not provided")
// 	// verify.That(endDate != nil, "Context is nil")

// 	// if err := verify.Flush(); err != nil {
// 	// 	return nil, err
// 	// }

// 	return b.query.GetTransactionsInDateRange(
// 		b.ctx,
// 		orm.GetTransactionsInDateRangeParams{StartDate: startDate, EndDate: endDate},
// 	)
// }

func (b *Budget) AddNewTransactionsFromDocument(fileName string, file io.Reader) error {
	log.Println("Processing transactions")

	verify := errors.Verifier{}

	verify.That(b.conn != nil, "Database connection is nil")
	verify.That(b.ctx != nil, "Context is nil")
	verify.That(fileName != "", "No Filename is provided")
	verify.That(file != nil, "File is nil")
	verify.That(
		strings.HasPrefix(strings.ToLower(fileName), "chase activity"),
		"Filename does not start with 'chase activity'",
	)

	if err := verify.Flush(); err != nil {
		return err
	}

	// ! This logic is not correct.
	// ! Year is being assumed to be the current year instead of the year the file was created.
	// ! That should be done client side
	parts := strings.Split(strings.ToLower(fileName), "activity")
	datePart := strings.TrimSpace(parts[1])
	datePart = strings.Split(datePart, ".")[0]
	datePart = strings.TrimSpace(datePart)

	dateParts := strings.Split(datePart, " ")
	if len(dateParts) != 2 {
		return errors.NewParseError(fmt.Sprintf("Error parsing date from filename. Date part: %v", datePart))
	}

	month, day := dateParts[0], dateParts[1]
	datePart = month + " " + day

	// Parse the extracted date
	parsedDate, err := time.Parse("Jan 2", datePart)
	if err != nil {
		return stdErrs.New("Error parsing date from filename: " + err.Error())
	}

	publishing_date := time.Date(time.Now().Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
	log.Println("Publishing date:", publishing_date)

	db := orm.New(b.conn)

	document_id, err := db.FindOneDocumentMeta(b.ctx, orm.FindOneDocumentMetaParams{
		Name:           fileName,
		PublishingDate: sql.NullTime{Time: publishing_date, Valid: true},
	})

	if err != nil && err != sql.ErrNoRows {
		return stdErrs.New("Error checking for existing document: " + err.Error())
	}

	if document_id == ""{
		println("we in the money")
	} else {
		println("we not in the money")
	}
	
	return nil
}
