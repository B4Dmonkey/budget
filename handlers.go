package main

import (
	"log"
	"my-budget/app"
	"net/http"
)

func root(ctx app.Context) error {
	log.Println("Root handler")
	homePage := HomePage{ctx: ctx}
	return ctx.Render(http.StatusOK, homePage)
}

func documents(ctx app.Context) error {
	log.Println("Documents handler")
	file, header, err := ctx.Req.FormFile("file")
	if err != nil {
		log.Println("Error reading file:", err)
		return err
	}
	defer file.Close()

	// * save the file to disk
	if err := SaveFileToDisk(header, file); err != nil {
		log.Println("Error saving file to upload dir:", err)
		return err
	}
	log.Println("File saved to disk")

	// * save the location to the database
	documentID, err := PersistDocumentMetaData(ctx.Req.Context(), header, file)
	if err != nil {
		log.Println("Error saving document metadata:", err)
		return err
	}
	log.Println("Document metadata saved")

	// * go through the db and create or update the records
	if err := PersistTransactions(ctx.Req.Context(), documentID, header, file); err != nil {
		log.Println("Error saving document metadata:", err)
		return err
	}
	log.Println("Transactions saved")

	// * respond with the values
	// Log the file name
	log.Printf("Uploaded file: %s", header.Filename)

	// You can add further processing of the file here

	ctx.Res.WriteHeader(http.StatusOK)
	if _, err := ctx.Res.Write([]byte("File uploaded successfully")); err != nil {
		return err
	}
	return nil
}
