package app

import "github.com/cbroglie/mustache"

type View interface {
	Mapping() interface{}
	Template() (*mustache.Template, error)
	Render() (string, error)
}
