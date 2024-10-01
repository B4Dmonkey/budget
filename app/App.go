package app

import (
	"net/http"
)

type HandlerFunc func(Context) error
type AppConfigFunc func(*AppConfig)
type Queries interface{}

type AppConfig struct {
	mux       *http.ServeMux
	server    *http.Server
	dbQueries interface{}
}

type App struct {
	AppConfig
}

func defaultAppConfig() AppConfig {
	return AppConfig{
		mux: http.NewServeMux(),
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
	a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Req: Request{r},
			Res: Response{w},
			DB:	a.dbQueries,
		}
		handler(ctx)
	})
}

func (a *App) Listen(addr string) error {
	a.server = &http.Server{
		Addr:    addr,
		Handler: a.mux,
	}
	return a.server.ListenAndServe()
}

func (a *App) Get(pattern string, handler HandlerFunc)  { a.AddHandler(MethodGet, pattern, handler) }
func (a *App) Post(pattern string, handler HandlerFunc) { a.AddHandler(MethodPost, pattern, handler) }
