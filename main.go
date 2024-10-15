package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sqlite "github.com/mattn/go-sqlite3"

	"my-budget/database/orm"
)

//go:embed database/schema.sql
var ddl string

//go:embed views/*
var viewsDir embed.FS

var CONNECT_ONCE sync.Once

const (
	Details int = iota
	PostingDate
	Description
	Amount
	Type
	Balance
)

func init() {
	log.Println("Extending sqlite3 driver...")
	CONNECT_ONCE.Do(func() {
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

type AppConfigFunc func(*AppConfig)

type AppConfig struct {
	ctx        context.Context
	db_conn    *sql.DB
	db_queries *orm.Queries
	mux        *http.ServeMux
	server     *http.Server
}

type App struct {
	AppConfig
}

func defaultAppConfig(ctx context.Context) (*AppConfig, error) {
	var ok bool

	var addr string

	if addr, ok = os.LookupEnv("DEVELOPMENT_PORT"); !ok || addr == "" {
		return nil, errors.New("DEVELOPMENT_PORT is not set")
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &AppConfig{ctx: ctx, mux: mux, server: server}, nil
}

func WithMux(mux *http.ServeMux) AppConfigFunc {
	return func(a *AppConfig) { a.mux = mux }
}

func WithServer(server *http.Server) AppConfigFunc {
	return func(a *AppConfig) { a.server = server }
}

func WithDbConnection(conn *sql.DB) AppConfigFunc {
	return func(a *AppConfig) { a.db_conn = conn }
}

func WithOrmQueries(queries *orm.Queries) AppConfigFunc {
	return func(a *AppConfig) { a.db_queries = queries }
}

func New(ctx context.Context, overrides ...AppConfigFunc) *App {
	var errs []error

	app_config, err := defaultAppConfig(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	for _, override := range overrides {
		override(app_config)
	}

	if app_config.db_conn == nil {
		conn, err := ConnectToDatabase(ctx)
		if err != nil {
			errs = append(errs, err)
		} else {
			app_config.db_conn = conn
			app_config.db_queries = orm.New(conn)
		}
	}

	app := App{
		AppConfig: *app_config,
	}

	if len(errs) > 0 {
		log.Fatalf("Failed to create app: %v", errs)
	}

	app.addHandlers()

	return &app
}

// func (a *App) Render(view View) (string, error) {
// ! The abstraction is too early
// 	if _, ok := view.(RenderOverridden); ok {
// 		return view.(RenderOverridden).Render()
// 	}

// 	template, err := view.Template()
// 	if err != nil {
// 		return "", err
// 	}
// 	viewBinding := view.Binding()
// 	return template.Render(viewBinding)
// }

func (a *App) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Println("Listen and server at", a.server.Addr)

		err := a.server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	<-a.ctx.Done()
	log.Println("Shutting down HTTP server gracefully...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelShutdown()

	err := a.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("HTTP server shutdown error: %s\n", err)
	}
}

func ConnectToDatabase(ctx context.Context) (*sql.DB, error) {
	log.Println("Connecting to database...")

	var connection *sql.DB

	var ok bool

	var err error

	var DATABASE_LOC string

	if DATABASE_LOC, ok = os.LookupEnv("DATABASE_LOC"); !ok || DATABASE_LOC == "" {
		return nil, errors.New("DATABASE_LOC is not set")
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	app := New(ctx)

	go app.Run(&wg)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()
	wg.Wait()
	log.Println("HTTP server stopped")
}
