package app

import "net/http"

type HandlerFunc func(Context) error

type Config struct {
	RequestMethods []string
}
type App struct {
	mux    *http.ServeMux
	server *http.Server
}

func New() *App {
	mux := http.NewServeMux()

	return &App{
		mux: mux,
	}
}

func (a *App) AddHandler(method string, path string, handler HandlerFunc) {
	// Todo: add defensive coding to restrict method and check pattern
	pattern := method + " " + path
	a.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := Context{
			Req: Request{r},
			Res: Response{w},
		}
		handler(ctx)
	})
}

func (a *App) Get(pattern string, handler HandlerFunc)  { a.AddHandler(MethodGet, pattern, handler) }
func (a *App) Post(pattern string, handler HandlerFunc) { a.AddHandler(MethodPost, pattern, handler) }
