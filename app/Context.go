package app

type Context struct {
	Req Request 
	Res Response
	Next HandlerFunc
}

// func (c *Context) Render(view View) error {
// 	return view.Render()
// }