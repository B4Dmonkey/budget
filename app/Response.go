package app

import "net/http"

type Response struct {
	w http.ResponseWriter
}

func (r *Response) Status(statusCode int)  error {
	r.w.WriteHeader(statusCode)
	return  nil
}