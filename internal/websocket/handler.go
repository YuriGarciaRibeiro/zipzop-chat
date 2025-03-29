package websocket

import (
	"log"
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, authService *auth.AuthService) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token não fornecido", http.StatusUnauthorized)
		return
	}

	email, err := authService.ValidateToken(token)
	if err != nil {
		log.Printf("Token inválido: %s | Erro: %v", token, err)
		http.Error(w, "Autenticação falhou", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erro no upgrade WebSocket: %v", err)
		return
	}

	client := &Client{
		ID:    uuid.New().String(),
		Email: email,
		Conn:  conn,
		Pool:  hub,
		Send:  make(chan Message, 256),
	}
	
	hub.register <- client

	go client.readPump()
	go client.writePump()

	log.Printf("Nova conexão estabelecida: %s (%s)", client.Email, client.ID)
}
