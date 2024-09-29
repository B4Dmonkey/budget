package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cbroglie/mustache"
)

type Context struct {
	responseWriter http.ResponseWriter
	request        *http.Request
}

func (ctx Context) Render(file_name_look_up string) error {
	log.Print("Rendering template")

	if file_name_look_up == "" {
		file_name_look_up = "views/pages/home.mst"
	}
	if file_name_look_up == "404page" {
		file_name_look_up = "views/pages/not-found.mst"
	}

	template, err := viewsDir.ReadFile(file_name_look_up)

	if err != nil {
		return errors.New("Error reading file: " + err.Error())
	}

	parsedTemplate, _ := mustache.ParseString(string(template))
	content, err := parsedTemplate.Render()
	if err != nil {
		return errors.New("Error rendering template: " + err.Error())
	}

	fmt.Fprint(ctx.responseWriter, content)

	return nil
}
