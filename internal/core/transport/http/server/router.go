package core_http_server

import (
	"fmt"
	"net/http"
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
}

func NewAPIVersionRoute(apiVersion APIVercion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   &http.ServeMux{},
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.Handler)
	}
}
