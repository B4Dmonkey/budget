package app

import (
	"net/http"
)

type HandlerFunc func(Context) error
type AppConfigFunc func(*AppConfig)
type Queries interface{}

type AppConfig struct {
	Mux       *http.ServeMux
	Server    *http.Server
	dbQueries interface{}
}

type App struct {
	AppConfig
}

func defaultAppConfig() AppConfig {
	return AppConfig{
		Mux: http.NewServeMux(),
	}
}

func WithDbQueries(dbQueries interface{}) AppConfigFunc {
	return func(cfg *AppConfig) {
		cfg.dbQueries = dbQueries
	}
}

func New(overrides ...AppConfigFunc) *App {
	cfg := defaultAppConfig()
	for _, override := range overrides {
		override(&cfg)
	}
	return &App{
		AppConfig: cfg,
	}
}

func (a *App) AddHandler(method string, path string, handler HandlerFunc) {
	// Todo: add defensive coding to restrict method and check pattern
	pattern := method + " " + path
	a.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Req: r,
			Res: w,
			DB:	a.dbQueries,
		}
		handler(ctx)
	})
}

func (a *App) Listen(addr string) error {
	a.Server = &http.Server{
		Addr:    addr,
		Handler: a.Mux,
	}
	return a.Server.ListenAndServe()
}

func (a *App) Get(pattern string, handler HandlerFunc)  { a.AddHandler(http.MethodGet, pattern, handler) }
func (a *App) Post(pattern string, handler HandlerFunc) { a.AddHandler(http.MethodPost, pattern, handler) }
