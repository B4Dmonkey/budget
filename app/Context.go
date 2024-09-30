package app

type Context struct {
	Req Request 
	Res Response
	Next HandlerFunc
}

func (c *Context) Render(statusCode int) error {
	return c.Res.Status(statusCode)
}
