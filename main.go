package main

import (
	// "context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	// "os"
	// "os/signal"
	// "sync"
	// "syscall"

	"github.com/cbroglie/mustache"
	sqlite "github.com/mattn/go-sqlite3"

	"my-budget/app"
	"my-budget/database/orm"
)

func init() {
	var connectOnce sync.Once
	log.Println("Extending sqlite3 driver...")
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
	log.Println("Extended sqlite3 driver successfully")
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

type View interface {
	Binding() interface{}
	Template() (*mustache.Template, error)
}

type RenderOverridden interface {
	Render() (string, error)
}

type App struct {
	mux        *http.ServeMux
	server     *http.Server
	db_conn    *sql.DB
	db_queries *orm.Queries
}

func ConnectDatabase() (*sql.DB, error) {
	var connection *sql.DB

	log.Println("Connecting to database...")
	var err error
	DATABASE_LOC := os.Getenv("DATABASE_LOC")

	if DATABASE_LOC == "" {
		return nil, fmt.Errorf("DATABASE_LOC is not set")
	}

	connection, err = sql.Open("sqlite3_extended", DATABASE_LOC)
	if connection == nil || err != nil {
		return nil, fmt.Errorf("Failed to connect the database: %s", DATABASE_LOC)
	}

	if connection.Ping() != nil {
		return nil, fmt.Errorf("Failed to ping database: %s", DATABASE_LOC)
	}

	log.Println("Database connected successfully")
	return connection, nil
}

func New() *App {
	var errs []error
	conn, err := ConnectDatabase()
	if err != nil {
		errs = append(errs, err)
	}

	mux := http.NewServeMux()
	app := App{
		db_conn: conn,
    db_queries: orm.New(conn),
		mux:     mux,
		server: &http.Server{
			Addr:    os.Getenv("DEVELOPMENT_PORT"),
			Handler: mux,
		},
	}

	if len(errs) > 0 {
		log.Fatalf("Failed to create app: %v", errs)
	}
	app.addHandlers()
	return &app
}

func New2() (*app.App, error) {
	if err := CreateDatabase(); err != nil {
		return nil, errors.New("Unable to create the file: " + err.Error())
	}

	app := app.New(
		app.WithDbQueries(orm.New(conn)),
	)

	// app.Get("/", root)
	app.Post("/documents", documents)

	return app, nil
}

func (a *App) Render(view View) (string, error) {
	if _, ok := view.(RenderOverridden); ok {
		return view.(RenderOverridden).Render()
	}

	template, err := view.Template()
	if err != nil {
		return "", err
	}
	viewBinding := view.Binding()
	return template.Render(viewBinding)
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// env := NewEnv()
	// var wg sync.WaitGroup
	// wg.Add(1)
	app := New()
	app.server.ListenAndServe()

	// if err := app.Listen(env.Addr()); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
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
