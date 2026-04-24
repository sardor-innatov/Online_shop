package route

import (
	"log"
	"net/http"
	"online_shop/src/common/config"
	"time"
)

type Router interface {
	Handle(method, path string, handler HandlerFunc)
	Group(prefix string, mw ...Middleware) *Group
	WrapHandler(h http.Handler) HandlerFunc
	Use(mw ...Middleware)
	Routes() []Route
	Start()
}

type router struct {
	mux         *http.ServeMux
	middlewares []Middleware
	routes      []Route
}

type Route struct {
	Path   string
	Method string
}

func NewRoute(mux *http.ServeMux) Router {
	return &router{
		mux:         mux,
		middlewares: make([]Middleware, 0),
		routes:      make([]Route, 0),
	}
}

type HandlerFunc func(Context) error

func makeHTTPHandlerFunc(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := Context{
			Response: w,
			Request:  r,
		}

		if err := h(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func makeHTTPHandler(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := Context{
			Response: w,
			Request:  r,
		}

		if err := h(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (r *router) Handle(method, path string, handler HandlerFunc) {

	fullpath := method + " " + path
	finalHandler := handler

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}

	r.mux.Handle(fullpath, makeHTTPHandlerFunc(finalHandler))

	newRoute := Route{
		Path:   path,
		Method: method,
	}
	r.routes = append(r.routes, newRoute)
}

func (r *router) Group(prefix string, mw ...Middleware) *Group {

	combinedMW := append(r.middlewares, mw...)

	return &Group{
		prefix:      prefix,
		middlewares: combinedMW,
		mux:         r.mux,
		router:      r,
	}
}

func (r *router) Start() {

	projectEnv := config.ProjectEnv()

	server := &http.Server{
		Addr:              projectEnv.ServerAddr,
		Handler:           r.mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second, // Защита от Slowloris
	}

	err := server.ListenAndServe()
	{
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}

}

func (r *router) WrapHandler(h http.Handler) HandlerFunc {
	return func(ctx Context) error {

		h.ServeHTTP(ctx.Response, ctx.Request)
		return nil
	}
}

func (r *router) Use(mw ...Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *router) Routes() []Route {
	return r.routes
}