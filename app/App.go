package app

import "net/http"

type HandlerFunc func(Context) error

type App struct {
	mux *http.ServeMux
}

func New() *App {
	mux := http.NewServeMux()

	return &App{
		mux: mux,
	}
}

func (a *App) AddHandler(method string, pattern string, handler HandlerFunc) {
	// Todo: add defensive coding to restrict method and check pattern
	pattern = method + " " + pattern
	a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Req: Request{r},
			Res: Response{w},
		}
		handler(ctx)
	})
}

func (a *App) Get(pattern string, handler HandlerFunc) { a.AddHandler(MethodGet, pattern, handler) }
func (a *App) Post(pattern string, handler HandlerFunc) { a.AddHandler(MethodPost, pattern, handler) }
