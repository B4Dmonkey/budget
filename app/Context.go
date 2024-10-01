package app

type Context struct {
	DB   interface{}
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

func (c *Context) Send(value string) error {
	if _, err := c.Res.w.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}
