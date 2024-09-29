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

func (a *App) Get(pattern string, handler HandlerFunc) {
	a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Req: Request{r},
			Res: Response{w},
		}
		handler(ctx)
	})
}
