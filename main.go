package main

import (
	"context"
	"database/sql"
	"embed"
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

type App struct {
	mux        *http.ServeMux
	server     *http.Server
	db_conn    *sql.DB
	db_queries *orm.Queries
}

func ConnectDatabase() (*sql.DB, error) {
	var connection *sql.DB

	log.Println("Connecting to database...")
	var ok bool
	var err error
	var DATABASE_LOC string

	if DATABASE_LOC, ok = os.LookupEnv("DATABASE_LOC"); !ok || DATABASE_LOC == "" {
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
		db_conn:    conn,
		db_queries: orm.New(conn),
		mux:        mux,
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

// func (a *App) InitializeServer() {
// Todo I need to grab the static stuff from this method. It's low priority
// 	if err := CreateDatabase(); err != nil {
// 		log.Fatal(err.Error())
// 		return
// 	}
// fSys, err := fs.Sub(publicDir, "public")

// if err != nil {
// 	log.Fatal("Failed to load public dir", err)
// }

// a.mux.Handle(GET+" /assets/", http.FileServer(http.FS(fSys)))

// a.mux.HandleFunc(GET+" /styles.css", func(w http.ResponseWriter, r *http.Request) {
// 	styles, err := publicDir.ReadFile("public/styles.css")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/css")
// 	w.Write(styles)
// })

// 	a.SetRouteHandlers(routes)
// }

// func (a *App) Render(view View) (string, error) {
// ! The abastraction is too early
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

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Println("Initializing HTTP server...")
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down HTTP server gracefully...")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err := a.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("HTTP server shutdown error: %s\n", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	app := New()

	go app.Run(ctx, &wg)
	app.server.ListenAndServe()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	cancel()
	wg.Wait()
	log.Println("HTTP server stopped")

}
