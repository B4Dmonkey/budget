package main

import "net/http"

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
