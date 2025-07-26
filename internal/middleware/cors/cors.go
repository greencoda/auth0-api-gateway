package cors

import (
	"github.com/gorilla/mux"
)

// ICORS interface defines the method to get the CORS middleware handler.
type ICORS interface {
	Handler() mux.MiddlewareFunc
}

// CORS implements the ICORS interface and provides the CORS middleware handler.
type CORS struct {
	middlewareFunc mux.MiddlewareFunc
}

// Handler returns the CORS middleware function.
func (c *CORS) Handler() mux.MiddlewareFunc {
	return c.middlewareFunc
}
