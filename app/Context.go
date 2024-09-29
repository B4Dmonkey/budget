package app

type Context struct {
	Req Request 
	Res Response
	Next HandlerFunc
}