package route

import (
	"net/http"
)

type Middleware func(HandlerFunc) HandlerFunc

type Group struct {
	prefix      string
	middlewares []Middleware
	mux         *http.ServeMux
	router      *router
}

func (g *Group) Handle(method, path string, handler HandlerFunc, mws ...Middleware) {

	fullpath := method + " " + g.prefix + path

	finalHandler := handler

	Middlewares := append(g.middlewares, mws...)

	for i := len(Middlewares) - 1; i >= 0; i-- { // reversing middliwares list
		finalHandler = Middlewares[i](finalHandler)
	}

	g.mux.Handle(fullpath, makeHTTPHandlerFunc(finalHandler))

	newRoute := Route{
		Path:   g.prefix + path,
		Method: method,
	}
	g.router.routes = append(g.router.routes, newRoute)
}

func (g *Group) Group(prefix string, mw ...Middleware) *Group {

	return &Group{
		prefix:      g.prefix + prefix,
		middlewares: append(g.middlewares, mw...), // inhereting middlewares and adding new
		mux:         g.mux,
		router:      g.router,
	}
}

func (g *Group) Use(mw ...Middleware) {
	g.middlewares = append(g.middlewares, mw...)
}
