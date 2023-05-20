package http

import "net/http"

type Middleware func(http.Handler) http.Handler

type MiddlewareMux struct {
	*http.ServeMux
	middleware []Middleware
}

func NewMiddlewareMux(middleware ...Middleware) *MiddlewareMux {
	return &MiddlewareMux{
		ServeMux:   http.NewServeMux(),
		middleware: middleware,
	}
}

func (m *MiddlewareMux) Use(middleware ...Middleware) {
	m.middleware = append(m.middleware, middleware...)
}

// ServeHTTP for MiddlewareMux ensures the sequence of middleware is preserved in order of Use calls
func (m *MiddlewareMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = m.ServeMux
	// preserve linear order, first .Use is first executed
	for i := range m.middleware {
		h = m.middleware[len(m.middleware)-1-i](h)
	}

	h.ServeHTTP(w, r)
}
