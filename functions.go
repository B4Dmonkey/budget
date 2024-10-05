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

func ConvertCurrencyIntToString(amount int64) string {
	return fmt.Sprintf("%.2f", float64(amount)/100)
}

func ConvertCurrencyStringToInt(amount string) (int64, error) {
	if amount == "" {
		log.Println("Empty amount returning nil")
		return 0, errors.New("Received empty string in ConvertCurrencyStringToInt")
	}
	amount = strings.Replace(amount, ".", "", -1)
	return strconv.ParseInt(amount, 10, 64)
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
    
    if err != nil {
      // ? Is this the best way to handle this ? I shouldn't run into this error
      if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err != csv.ErrFieldCount {
        return err
      }
    }

		amount, err := ConvertCurrencyStringToInt(record[Amount])
		if err != nil {
			log.Println("Error parsing amount:", err)
			continue
		}

		var balance sql.NullInt64
		if strings.TrimSpace(record[Balance]) == "" {
			balance = sql.NullInt64{Valid: false}
		} else {
			parsed_balance, err := ConvertCurrencyStringToInt(record[Balance])
			if err != nil {
				log.Println("Error parsing balance:", err)
				continue
			}
			balance = sql.NullInt64{Int64: parsed_balance, Valid: true}
		}

		postingDate, err := time.Parse("01/02/2006", record[PostingDate])
		if err != nil {
			log.Println("Error parsing posting date:", err)
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
