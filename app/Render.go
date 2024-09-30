package app

type View interface {
	Render() (string, error)
	Template() string
}
