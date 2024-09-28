package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"my-budget/database/orm"
	"os"
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

const (
	Details int = iota
	PostingDate
	Description
	Amount
	Type
	Balance
)

func PersistTransactions(ctx context.Context, header *multipart.FileHeader, file multipart.File) error {
	// tx, err := conn.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer tx.Rollback()

	// db := orm.New(conn)
	// txdb := db.WithTx(tx)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("error skipping headers: %v", err)
	}

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

		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err == csv.ErrFieldCount {
			log.Println("Error parsing record:", record)
		}

		log.Println("CSV Record:", record)
	}

	//

	// transactions, err := db.CreateTransactions(ctx, orm.CreateTransactionsParams{
	// 	Name: header.Filename,
	// })

	// if err != nil {
	// 	return errors.New("Error creating transactions record: " + err.Error())
	// }

	// log.Println("Transactions created:", transactions.ID)

	return nil
}
