package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/DimaKirejko/todoapp/internal/core/transport/http/middleware"
)

type APIVercion string

var (
	ApiVersion1 = APIVercion("v1")
	ApiVersion2 = APIVercion("v2")
	ApiVersion3 = APIVercion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVercion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRoute(apiVersion APIVercion, middleware ...core_http_middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   &http.ServeMux{},
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
