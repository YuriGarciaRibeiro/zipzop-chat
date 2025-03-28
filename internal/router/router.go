package router

import (
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/gorilla/mux"
)

type Route struct {
	URI      string
	Method   string
	Function func(w http.ResponseWriter, r *http.Request)
	IsAuth   bool
}

func ConfigureAuthRoutes(r *mux.Router, authHandler *auth.AuthHandler) {
	routes := NewAuthRoutes(authHandler)
	for _, route := range routes {
		if route.IsAuth {
			// TODO: Implementar middleware de autenticação
			r.HandleFunc(route.URI, route.Function).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, route.Function).Methods(route.Method)
		}
	}
}
