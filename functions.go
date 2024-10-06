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

	"my-budget/database/orm"
)

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

func ProcessTransactions(db orm.Querier, ctx context.Context, header *multipart.FileHeader, file io.Reader) error {
	log.Println("Processing transactions")
	document_id, err := db.FindOneDocumentMeta(ctx, header.Filename)

	if err != nil && err != sql.ErrNoRows {
		return errors.New("Error checking for existing document: " + err.Error())
	}

	if document_id != "" {
		log.Println("Document already exists")
		// 	return ProcessNewTransactions(ctx, document_id, file)
	} else {
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

func PersistDocumentMetaData(ctx context.Context, header *multipart.FileHeader, file multipart.File) (string, error) {
	db := orm.New(conn)

	document_id, err := db.FindOneDocumentMeta(ctx, header.Filename)

	if err != nil && err != sql.ErrNoRows {
		return "", errors.New("Error checking for existing document: " + err.Error())
	}

	if document_id != "" {
		return document_id, nil
	}

	document, err := db.CreateDocumentMeta(ctx, orm.CreateDocumentMetaParams{
		Name:         header.Filename,
		PersistedLoc: header.Filename,
	})

	if err != nil {
		return "", errors.New("Error creating document record: " + err.Error())
	}

	log.Println("Document created:", document.ID)

	return document.ID, nil
}

/*
* Get all transactions which are pending
* loop through csv
* check if the transaction match an elem in the array
* if it does, update the transaction
* if it doesn't, create a new transaction
* end loop when i get to the end of pending transactions
 */

func PersistTransactions(ctx context.Context, documentID string, header *multipart.FileHeader, file multipart.File) error {
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	db := orm.New(conn)
	txdb := db.WithTx(tx)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	pending_transactions, err := db.GetPendingTransactions(ctx, documentID)
	if err != nil {
		return err
	}

	log.Println("Pending transactions:", pending_transactions)

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error skipping headers: %v", err)
	}

	log.Println("Looping through records...")
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

		var isMatch bool
		for _, transaction := range pending_transactions {
			details_match := transaction.Details == record[Details]
			postingDate_match := transaction.PostingDate == postingDate
			description_match := transaction.Description == record[Description]
			transactionType_match := transaction.Type == record[Type]
			amount_match := transaction.Amount == amount
			balance_match := transaction.Balance == balance

			if details_match && postingDate_match && description_match && amount_match && transactionType_match && balance_match {
				isMatch = true
				break
			}
			isMatch = false
		}

		if isMatch {
			continue
		}

		if err := txdb.CreateTransaction(ctx, orm.CreateTransactionParams{
			DocumentID:  documentID,
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
	tx.Commit()
	log.Println("Transactions created")
	return nil
}

func SaveFileToDisk(header *multipart.FileHeader, file multipart.File) error {
	// todo: check if the folder is there. Not important since this is for me
	out, err := os.Create("uploads/" + header.Filename)
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
