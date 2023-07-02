package middlewares

import (
	"net/http"
)

type MiddlewareHandler struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func NewMiddlewareHandler() *MiddlewareHandler {
	return &MiddlewareHandler{
		mux: http.NewServeMux(),
	}
}

func (mh *MiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//direct the request to the embedded servemux
	mh.mux.ServeHTTP(w, r)
}

func (mh *MiddlewareHandler) Handle(pattern string, handler http.Handler) {
	// Apply the middlewares to the handler
	for _, middleware := range mh.middlewares {
		handler = middleware(handler)
	}

	// Register the handler with the underlying ServeMux
	mh.mux.Handle(pattern, handler)
}

func (mh *MiddlewareHandler) Use(middleware func(http.Handler) http.Handler) {
	// Append the middleware to the middlewares slice
	mh.middlewares = append(mh.middlewares, middleware)
}
