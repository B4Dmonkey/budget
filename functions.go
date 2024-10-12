package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"my-budget/database/orm"
)

const DocumentUploads = "DocumentUploads"

func currentTimestamp() string { return time.Now().Format("2006-01-02 15:04:05") }

func newUUID() string { return uuid.New().String() }

func ConvertCurrencyIntToString(amount int64) string { return fmt.Sprintf("%.2f", float64(amount)/100) }

func ConvertCurrencyStringToInt64(amount string) (int64, error) {
	return strconv.ParseInt(strings.Replace(strings.TrimSpace(amount), ".", "", -1), 10, 64)
}

func ConvertCurrencyStringToNullInt64(amount string) sql.NullInt64 {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return sql.NullInt64{Valid: false}
	}

	int64_amount, err := ConvertCurrencyStringToInt64(amount)

	if err != nil {
		// Todo: I should check what errors this might be and handle them
		log.Println("Error converting currency string to int:", err)
	}

	return sql.NullInt64{
		Int64: int64_amount,
		Valid: err == nil,
	}
}

func GetTransactions(db orm.Querier, ctx context.Context, startDate time.Time, endDate time.Time) ([]orm.Transaction, error) {
	verify := Verifier{}
	verify.That(db != nil, "Database connection is nil")
	verify.That(ctx != nil, "Context is nil")
	if err := verify.Flush(); err != nil {
		return nil, err
	}

	return db.GetTransactionsInDateRange(
		ctx,
		orm.GetTransactionsInDateRangeParams{PostingDate: startDate, PostingDate_2: endDate},
	)
}

func ProcessNewTransactions(db orm.Querier, ctx context.Context, header *multipart.FileHeader, file io.Reader) error {
	log.Println("Processing transactions")
	verify := Verifier{}
	verify.That(db != nil, "Database connection is nil")
	verify.That(ctx != nil, "Context is nil")
	verify.That(header != nil, "Header is nil")
	verify.That(file != nil, "File is nil")
	if err := verify.Flush(); err != nil {
		return err
	}

	fileName := strings.ToLower(header.Filename)
	if !strings.HasPrefix(fileName, "chase activity") {
		return &VerificationError{Message: fmt.Sprintf("Filename '%s' does not start with 'chase activity'", header.Filename)}
	}

	if err := SaveFileToDisk(header, file); err != nil {
		return err
	}

	// ! This logic is not correct. Year is being assumed to be the current year instead of the year the file was created. That should be done client side
	parts := strings.Split(fileName, "activity")
	datePart := strings.TrimSpace(parts[1])
	datePart = strings.Split(datePart, ".")[0]
	datePart = strings.TrimSpace(datePart)

	monthMap := map[string]string{
		"sept": "Sep", // Handle "Sept" specifically
	}

	dateParts := strings.Split(datePart, " ")
	if len(dateParts) != 2 {
		return newParseError(fmt.Sprintf("Error parsing date from filename. Date part: %v", datePart))
	}

	month, day := dateParts[0], dateParts[1]
	if abbr, ok := monthMap[month]; ok {
		month = abbr
	} else {
		return errors.New("Error parsing month from filename")
	}

	datePart = month + " " + day

	// Parse the extracted date
	parsedDate, err := time.Parse("Jan 2", datePart)
	if err != nil {
		return errors.New("Error parsing date from filename: " + err.Error())
	}

	// Combine with the year from creationTime
	publishing_date := time.Date(time.Now().Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
	log.Println("Publishing date:", publishing_date)

	document_id, err := db.FindOneDocumentMeta(ctx, orm.FindOneDocumentMetaParams{
		Name:           header.Filename,
		PublishingDate: sql.NullTime{Time: publishing_date, Valid: true},
	})

	if err != nil && err != sql.ErrNoRows {
		return errors.New("Error checking for existing document: " + err.Error())
	}

	//todo: maybe send message for duplicate documents
	if document_id == "" {
		log.Println("Creating new document")
		document, err := db.CreateDocumentMeta(ctx, orm.CreateDocumentMetaParams{
			Name:         header.Filename,
			PersistedLoc: header.Filename,
		})

		if err != nil {
			return errors.New("Error creating document record: " + err.Error())
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

			amount, err := ConvertCurrencyStringToInt64(record[Amount])
			if err != nil {
				log.Println("Error parsing amount:", err)
				continue
			}

			balance := ConvertCurrencyStringToNullInt64(record[Balance])

			postingDate, err := time.Parse("01/02/2006", record[PostingDate])
			if err != nil {
				log.Println("Error parsing posting date:", err)
				continue
			}

			if err := db.CreateTransaction(ctx, orm.CreateTransactionParams{
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
	}

	return nil
}

func SaveFileToDisk(header *multipart.FileHeader, file io.Reader) error {
	// todo: check if the folder is there. Not important since this is for me
	out, err := os.Create(filepath.Join(DocumentUploads, header.Filename))
	if err != nil {
		return errors.New("Unable to create the file: " + err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return errors.New("Error saving file: " + err.Error())
	}

	return nil
}
