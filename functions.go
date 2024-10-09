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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"my-budget/database/orm"
)

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

func GetTransactions(db orm.Querier, ctx context.Context) ([]orm.Transaction, error) {
	if db == nil {
		return nil, errors.New("Database connection is nil")
	}
	if ctx == nil {
		return nil, errors.New("Context is nil")
	}
	
	return nil, nil
}

func ProcessNewTransactions(db orm.Querier, ctx context.Context, header *multipart.FileHeader, file io.Reader) error {
	log.Println("Processing transactions")
	if db == nil {
		return errors.New("Database connection is nil")
	}
	if ctx == nil {
		return errors.New("Context is nil")
	}
	if header == nil {
		return errors.New("Header is nil")
	}
	if file == nil {
		return errors.New("File is nil")
	}
	document_id, err := db.FindOneDocumentMeta(ctx, header.Filename)

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

func SaveFileToDisk(header *multipart.FileHeader, file multipart.File) error {
	// todo: check if the folder is there. Not important since this is for me
	out, err := os.Create("DocumentUploads/" + header.Filename)
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
