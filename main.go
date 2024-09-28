package main

import (
	"context"
	"embed"
	"fmt"

	// "io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cbroglie/mustache"
)

// //go:embed public/*
// var publicDir embed.FS

//go:embed views/*
var viewsDir embed.FS

type EnvConfigInterface interface {
	IsDev() bool
	Addr() string
}

type EnvConfig struct {
	isDev bool
	addr  string
}

func NewEnv() *EnvConfig {
	log.Println("Loading environment variables...")

	var addr string
	addr = os.Getenv("DEVELOPMENT_PORT")
	return &EnvConfig{ addr: addr}
}

func (e *EnvConfig) IsDev() bool {
	return e.isDev
}

func (e *EnvConfig) Addr() string {
	return e.addr
}

type Route struct {
	httpVerb string
	pattern  string
	handler  http.HandlerFunc
}

var routes = []Route{
	{
		GET,
		"/",
		WithMiddleware(root, logger),
	},
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func WithMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// * loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

type App struct {
	mux    *http.ServeMux
	server *http.Server
	ctx    context.Context
}

func CreateApp(ctx context.Context, env EnvConfigInterface) *App {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    env.Addr(),
		Handler: mux,
	}

	return &App{
		mux:    mux,
		server: server,
		ctx:    ctx,
	}
}

const (
	GET  = "GET"
	POST = "POST"
)

func (a *App) InitializeServer() {
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

	a.SetRouteHandlers(routes)
}

func (a *App) SetRouteHandlers(handlers []Route) {
	log.Println("Adding routes...")
	for _, r := range handlers {
		route := r.httpVerb + " " + r.pattern
		// Todo: this should be a -v log
		// log.Println(route)
		a.mux.HandleFunc(route, r.handler)
	}
	log.Println("Routes added")
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		log.Println("Initializing HTTP server...")
		a.InitializeServer()
		log.Println("Starting HTTP server...")
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

func Render(w http.ResponseWriter, file_name_look_up string) {
	log.Print("Rendering template")
	if file_name_look_up == "" {
		file_name_look_up = "views/pages/home.mst"
	}
	if file_name_look_up == "404page" {
		file_name_look_up = "views/pages/not-found.mst"
	}

	template, err := viewsDir.ReadFile(file_name_look_up)
	log.Println(file_name_look_up)
	if err != nil {
		log.Println("Error reading file", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	parsedTemplate, _ := mustache.ParseString(string(template))
	content, err := parsedTemplate.Render()
	if err != nil {
		log.Println("Error rendering template", err)
		http.Error(w, "Internal server error two", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}

func root(w http.ResponseWriter, r *http.Request) {
	log.Println("root handler")
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		Render(w, "404page")
		return
	}
	Render(w, "")
}

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
