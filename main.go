package main

import (
	// "context"
	"database/sql"
	"embed"
	"errors"
	"log"
	"sync"

	// "os"
	// "os/signal"
	// "sync"
	// "syscall"

	sqlite "github.com/mattn/go-sqlite3"

	"my-budget/app"
	"my-budget/database/orm"
)

func init() {
	var connectOnce sync.Once
	connectOnce.Do(func() {
		sql.Register("sqlite3_extended", &sqlite.SQLiteDriver{
			ConnectHook: func(conn *sqlite.SQLiteConn) error {

				if err := conn.RegisterFunc("uuid", newUUID, false); err != nil {
					return err
				}

				if err := conn.RegisterFunc("current_timestamp", currentTimestamp, false); err != nil {
					return err
				}

				return nil
			},
		})
	})
}

// //go:embed public/*
// var publicDir embed.FS

//go:embed views/*
var viewsDir embed.FS

const (
	Details int = iota
	PostingDate
	Description
	Amount
	Type
	Balance
)

func New() (*app.App, error) {
	if err := CreateDatabase(); err != nil {
		return nil, errors.New("Unable to create the file: " + err.Error())
	}

	app := app.New(
		app.WithDbQueries(orm.New(conn)),
	)

	app.Get("/", root)
	app.Post("/documents", documents)

	return app, nil
}

func main() {
	app, err := New()
	env := NewEnv()

	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := app.Listen(env.Addr()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// func main2() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	env := NewEnv()
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	app := CreateApp(ctx, env)
// 	go app.Start(&wg)

// 	signalCh := make(chan os.Signal, 1)
// 	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

// 	<-signalCh
// 	cancel()
// 	wg.Wait()
// 	log.Println("HTTP server stopped")
// }
