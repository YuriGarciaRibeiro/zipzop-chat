package server

import (
	"context"
	"log"
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/repository"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

func NewServer(cfg *config.AppConfig) *Server {
	router := mux.NewRouter()

	return &Server{
		router: router,
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
	userRepo *repository.UserRepository,
) {
	authHandler := auth.NewAuthHandler(authService)
	authMiddleware := auth.NewAuthMiddleware(authService)

	// Rotas públicas
	s.router.HandleFunc("/health", healthCheck).Methods("GET")
	s.router.HandleFunc("/register", authHandler.Register).Methods("POST")
	s.router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Rotas protegidas
	protected := s.router.PathPrefix("/api").Subrouter()
	protected.Use(authMiddleware.Handler)

	// Exemplo de rota protegida
	protected.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		// Acesse o userID do contexto
		userID := r.Context().Value("userID").(string)
		w.Write([]byte("User ID: " + userID))
	}).Methods("GET")
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
