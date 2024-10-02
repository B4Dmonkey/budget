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
	"my-budget/database/orm"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func PersistDocumentMetaData(ctx context.Context, header *multipart.FileHeader, file multipart.File) (string, error) {
	db := orm.New(conn)

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
		if err != nil {
			if err == io.EOF {
				break
			}
			if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err != csv.ErrFieldCount {
				return err
			}
		}

		amount, err := strconv.ParseInt(strings.Replace(record[Amount], ".", "", -1), 10, 64)
		if err != nil {
			log.Println("Error parsing amount:", err)
			continue
		}

		var balance sql.NullInt64
		if strings.TrimSpace(record[Balance]) == "" {
			balance = sql.NullInt64{Valid: false}
		} else {
			parsed_balance, err := strconv.ParseInt(strings.Replace(record[Balance], ".", "", -1), 10, 64)
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
