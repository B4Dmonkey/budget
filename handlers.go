package main

import (
	"log"
	"net/http"

	"my-budget/app"
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

	if err := SaveFileToDisk(header, file); err != nil {
		log.Println("Error saving file to upload dir:", err)
		return err
	}
	log.Println("File saved to disk")

	documentID, err := PersistDocumentMetaData(ctx.Req.Context(), header, file)
	if err != nil {
		log.Println("Error saving document metadata:", err)
		return err
	}
	log.Println("Document metadata saved")

	if err := PersistTransactions(ctx.Req.Context(), documentID, header, file); err != nil {
		log.Println("Error saving document metadata:", err)
		return err
	}
	log.Println("Transactions saved")

	ctx.Res.WriteHeader(http.StatusOK)
	if _, err := ctx.Res.Write([]byte("File uploaded successfully")); err != nil {
		return err
	}
	return nil
}
