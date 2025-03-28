package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/router"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/websocket"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	hub        *websocket.Hub
}

func NewServer(cfg *config.AppConfig) *Server {
	router := mux.NewRouter()
	hub := websocket.NewHub()

	// Inicia o gerenciador de WebSocket em uma goroutine separada
	go hub.Start()

	return &Server{
		router: router,
		hub:    hub,
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) SetupRoutes(
	authService *auth.AuthService,
) {
	authHandler := auth.NewAuthHandler(authService)

	// Rotas públicas
	s.router.HandleFunc("/health", healthCheck).Methods("GET")

	router.ConfigureAuthRoutes(s.router, authHandler)

	// Websocket (definir ANTES do PathPrefix)
	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(s.hub, w, r, authService)
	})

	// Rotas estáticas (deve ser a última)
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	fmt.Println("Rotas configuradas corretamente")
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
