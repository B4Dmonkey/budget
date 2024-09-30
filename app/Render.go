package app

import "github.com/cbroglie/mustache"

type View interface {
	Mapping() interface{}
	Template() (*mustache.Template, error)
}

type RenderOverridden interface {
	Render() (string, error)
}

func Render(view View) (string, error) {
	if _, ok := view.(RenderOverridden); ok {
		return view.(RenderOverridden).Render()
	}

	template, err := view.Template()
	if err != nil {
		return "", err
	}
	mapping := view.Mapping()
	return template.Render(mapping)
}
