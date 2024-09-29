package main

import "net/http"

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}
