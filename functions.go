package main

import (
	"context"
	"errors"
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

func PersistDocumentMetaData(ctx context.Context, header *multipart.FileHeader, file multipart.File) ( error) {
	db := orm.New(conn)

	document, err := db.CreateDocument(ctx, orm.CreateDocumentParams{
		Name:         header.Filename,
		PersistedLoc: header.Filename,
	})

	if err != nil {
		return errors.New("Error creating document record: " + err.Error())
	}

	log.Println("Document created:", document.ID)

	return nil
}
