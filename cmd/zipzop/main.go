package main

import (
	"fmt"
	"log"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/config"
	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/server"
)

func main() {
	fmt.Println("Iniciando servidor...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Servidor iniciado na porta 8080")
	if err := server.StartServer(cfg); err != nil {
		log.Fatal(err)
	}
}
