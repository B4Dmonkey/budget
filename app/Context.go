package app

type Context struct {
	Req  Request
	Res  Response
	Next HandlerFunc
}

func (c *Context) Render(statusCode int, view View) error {
	if err := c.Res.Status(statusCode); err != nil {
		return err
	}
	content, err := Render(view)
	if err != nil {
		return err
	}
	_, err = c.Res.w.Write([]byte(content))
	return err
}
