package app

import "github.com/cbroglie/mustache"

type View interface {
	Binding() interface{}
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

	viewBinding := view.Binding()

	return template.Render(viewBinding)
}
