package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/repository"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/server"
)

func main() {
	// Carrega configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inicializa repositórios
	userRepo, err := repository.NewUserRepository(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// Inicializa serviços
	authService := auth.NewAuthService(userRepo, cfg)

	// Configura servidor HTTP
	srv := server.NewServer(cfg)
	srv.SetupRoutes(authService, userRepo)

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	log.Println("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
