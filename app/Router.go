package app

type Router interface {
	Use(args ...any) Router

	GET(path string, handler HandlerFunc, middleware ...HandlerFunc) Router
	POST(path string, handler HandlerFunc, middleware ...HandlerFunc) Router

	AddHandler(method string, path string, handler HandlerFunc) Router
	All(path string, handler HandlerFunc, middleware ...HandlerFunc) Router
}