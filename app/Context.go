package app

import "net/http"

type Context struct {
	DB   interface{}
	Req  *http.Request
	Res  http.ResponseWriter
	Next HandlerFunc
}

func (c *Context) Render(statusCode int, view View) error {
	c.Status(statusCode)
	content, err := Render(view)
	if err != nil {
		return err
	}
	_, err = c.Res.Write([]byte(content))
	return err
}

func (c *Context) Send(value string) error {
	if _, err := c.Res.Write([]byte(value)); err != nil {
		return err
	}
	return nil
}

func (r *Context) Status(statusCode int) error {
	r.Res.WriteHeader(statusCode)
	return nil
}
