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

	if err := ProcessNewTransactions(txdb, ctx.Req.Context(), header, file); err != nil {
		log.Println("Error processing transactions:", err)
		return err
	}

	tx.Commit()
	return nil
}

