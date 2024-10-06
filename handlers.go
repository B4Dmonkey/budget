package main

import (
	"log"
	"net/http"

	"my-budget/app"
	"my-budget/database/orm"
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

	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	db := orm.New(conn)
	txdb := db.WithTx(tx)

	if err := ProcessTransactions(txdb, ctx.Req.Context(), header, file); err != nil {
		log.Println("Error processing transactions:", err)
		return err
	}

	tx.Commit()
	return nil
}

func documents2(ctx app.Context) error {
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
