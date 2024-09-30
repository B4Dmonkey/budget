package app

import "github.com/cbroglie/mustache"

type View interface {
	Mapping() interface{}
	Template() (*mustache.Template, error)
}

func Render(view View) (string, error) {
	template, err := view.Template()
	if err != nil {
		return "", err
	}
	mapping := view.Mapping()
	return template.Render(mapping)
}
