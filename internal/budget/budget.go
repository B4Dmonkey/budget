package budget

import (
	"context"
	"database/sql"
	"encoding/csv"
	stdErrs "errors"
	"fmt"
	"io"
	"log"
	"my-budget/database/orm"
	"my-budget/internal/convert"
	"my-budget/internal/errors"
	"strings"
	"time"
)

// Constants used to index the columns in a transactions csv file
const (
	Details int = iota
	PostingDate
	Description
	Amount
	Type
	Balance
)

type Budget struct {
	ctx  context.Context
	conn *sql.DB
}

// Create a new Budget instance
// Errors if the context or database connection is nil
func NewBudget(ctx context.Context, conn *sql.DB) (*Budget, error) {
	verify := errors.Verifier{}
	verify.That(conn != nil, "Querier not provided")
	verify.That(ctx != nil, "Context is nil")

	if err := verify.Flush(); err != nil {
		return nil, err
	}

	return &Budget{ctx: ctx, conn: conn}, nil
}

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
	monthMap := map[string]string{
		"sept": "Sep", // Handle "Sept" specifically
	}

	if abbr, ok := monthMap[month]; ok {
		month = abbr
	} else {
		log.Println("Month not found in map, using original parsed name:", month)
	}
	datePart = month + " " + day

	// Parse the extracted date
	parsedDate, err := time.Parse("Jan 2", datePart)
	if err != nil {
		return stdErrs.New("Error parsing date from filename: " + err.Error())
	}

	publishing_date := time.Date(time.Now().Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
	log.Println("Publishing date:", publishing_date)

	db := orm.New(b.conn) // ? Is there a better way to do this ?

	document_id, err := db.FindOneDocumentMeta(b.ctx, orm.FindOneDocumentMetaParams{
		Name:           fileName,
		PublishingDate: sql.NullTime{Time: publishing_date, Valid: true},
	})

	if err != nil && err != sql.ErrNoRows {
		return stdErrs.New("Error checking for existing document: " + err.Error())
	}

	if document_id != "" {
		return stdErrs.New("Document Exist") // ? Document already exists. What should be done here?
	}

	log.Println("Document does not exist, creating new document")

	document, err := db.CreateDocumentMeta(b.ctx, orm.CreateDocumentMetaParams{
		Name:           fileName,
		PersistedLoc:   fileName,
		PublishingDate: sql.NullTime{Time: publishing_date, Valid: true},
	})

	if err != nil {
		return stdErrs.New("Error creating document record: " + err.Error())
	}

	// Reset the cursor to the start of the file
	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, 0) // Reset the cursor to the start of the file
		if err != nil {
			return fmt.Errorf("error seeking to beginning of file: %v", err)
		}
	}

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error skipping headers: %v", err)
	}

	log.Println("Document created:", document.ID)

	for {
		record, err := reader.Read()
		if err != nil && err == io.EOF {
			break
		}

		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err != csv.ErrFieldCount {
			return err
		}

		amount, err := convert.CurrencyStringToInt64(record[Amount])
		if err != nil {
			log.Println("Error parsing amount:", err)
			continue
		}

		balance := convert.CurrencyStringToNullInt64(record[Balance])

		postingDate, err := time.Parse("01/02/2006", record[PostingDate])
		if err != nil {
			log.Println("Error parsing posting date:", err)
			continue
		}

		if err := db.CreateTransaction(b.ctx, orm.CreateTransactionParams{
			DocumentID:  document.ID,
			Details:     record[Details],
			PostingDate: postingDate,
			Description: record[Description],
			Amount:      amount,
			Type:        record[Type],
			Balance:     balance,
		}); err != nil {
			log.Println("Error creating transaction record:", err)
			continue
		}
	}

	return nil
}

func (b *Budget) GetIncomeVsExpense() error {
	db := orm.New(b.conn)

	result, err := db.GetIncomeVsExpenseAndTotal(b.ctx)
	if err != nil {
		return stdErrs.New("Error getting income vs expense: " + err.Error())
	}

	log.Printf("Income: %v, Expense: %v, Total: %v\n", result.Income, result.Expense, result.Total)

	return nil
}
