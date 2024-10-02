package main

import (
	// "context"
	"embed"
	"errors"
	"my-budget/app"
	"my-budget/database/orm"

	"log"
	// "os"
	// "os/signal"
	// "sync"
	// "syscall"
)

// //go:embed public/*
// var publicDir embed.FS

//go:embed views/*
var viewsDir embed.FS

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

func New() (*app.App, error) {
	if err := CreateDatabase(); err != nil {
		return nil, errors.New("Unable to create the file: " + err.Error())
	}

	app := app.New(
		app.WithDbQueries(orm.New(conn)),
	)

	app.Get("/", root)
	app.Post("/upload", upload)

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
