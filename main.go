package main

import (
	"context"
	"embed"

	// "io/fs"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// //go:embed public/*
// var publicDir embed.FS

//go:embed views/*
var viewsDir embed.FS

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := NewEnv()
	var wg sync.WaitGroup
	wg.Add(1)
	app := CreateApp(ctx, env)
	go app.Start(&wg)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()
	wg.Wait()
	log.Println("HTTP server stopped")
}
