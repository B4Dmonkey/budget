package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"
//
//
// )

// type App struct {
// 	mux    *http.ServeMux
// 	server *http.Server
// 	ctx    context.Context
// }

// func CreateApp(ctx context.Context, env EnvConfigInterface) *App {
// 	mux := http.NewServeMux()
// 	server := &http.Server{
// 		Addr:    env.Addr(),
// 		Handler: mux,
// 	}

// 	return &App{
// 		mux:    mux,
// 		server: server,
// 		ctx:    ctx,
// 	}
// }

// func (a *App) InitializeServer() {
// 	if err := CreateDatabase(); err != nil {
// 		log.Fatal(err.Error())
// 		return
// 	}
// 	// fSys, err := fs.Sub(publicDir, "public")

// 	// if err != nil {
// 	// 	log.Fatal("Failed to load public dir", err)
// 	// }

// 	// a.mux.Handle(GET+" /assets/", http.FileServer(http.FS(fSys)))

// 	// a.mux.HandleFunc(GET+" /styles.css", func(w http.ResponseWriter, r *http.Request) {
// 	// 	styles, err := publicDir.ReadFile("public/styles.css")
// 	// 	if err != nil {
// 	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 	// 		return
// 	// 	}
// 	// 	w.Header().Set("Content-Type", "text/css")
// 	// 	w.Write(styles)
// 	// })

// 	a.SetRouteHandlers(routes)
// }

// func (a *App) SetRouteHandlers(handlers []Route) {
// 	log.Println("Adding routes...")
// 	for _, r := range handlers {
// 		route := r.httpVerb + " " + r.pattern
// 		// Todo: this should be a -v log
// 		// log.Println(route)
// 		a.mux.HandleFunc(route, r.handler)
// 	}
// 	log.Println("Routes added")
// }

// func (a *App) Start(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	go func() {
// 		log.Println("Initializing HTTP server...")
// 		a.InitializeServer()
// 		log.Println("Starting HTTP server...")
// 		err := a.server.ListenAndServe()
// 		if err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %s\n", err)
// 		}
// 	}()

// 	<-a.ctx.Done()
// 	log.Println("Shutting down HTTP server gracefully...")
// 	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelShutdown()

// 	err := a.server.Shutdown(shutdownCtx)
// 	if err != nil {
// 		log.Printf("HTTP server shutdown error: %s\n", err)
// 	}
// }
