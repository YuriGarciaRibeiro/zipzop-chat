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

type Router struct {
	r          *mux.Router
	middleware *auth.AuthMiddleware
}

func NewRouter(muxRouter *mux.Router, authMiddleware *auth.AuthMiddleware) *Router {
	return &Router{
		r:          muxRouter,
		middleware: authMiddleware,
	}
}

func (router *Router) configureRoutes(routes []Route) {
	authSubRouter := router.r.PathPrefix("/").Subrouter() // Subrouter para rotas autenticadas
	authSubRouter.Use(router.middleware.Handler)          // Aplica o middleware

	for _, route := range routes {
		if route.IsAuth {
			authSubRouter.HandleFunc(route.URI, route.Function).Methods(route.Method)
		} else {
			router.r.HandleFunc(route.URI, route.Function).Methods(route.Method)
		}
	}
}

func (router *Router) ConfigureAuthRoutes(authHandler *auth.AuthHandler) {
	routes := NewAuthRoutes(authHandler)
	router.configureRoutes(routes)
}
