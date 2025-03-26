package websocket

import (
	"log"
	"net/http"

	"github.com/YuriGarciaRibeiro/zipzop-chat/internal/auth"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWs gerencia conexões WebSocket e usa o userID do contexto
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, authService *auth.AuthService) {
	// Extrair o token da query string
	token := r.URL.Query().Get("token")

	if token == "" {
		log.Println("Token não fornecido")
		http.Error(w, "Token necessário", http.StatusUnauthorized)
		return
	}

	// Validar o token e obter o userID
	email, err := authService.ValidateToken(token)
	if err != nil {
		log.Println("Token inválido:", err)
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Conectar ao WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade para WebSocket:", err)
		return
	}

	// Criar um novo cliente WebSocket com o userID
	client := &Client{
		Conn:  conn,
		Pool:  hub,
		email: email, 
		Send:  make(chan Message),
	}

	// Registrar o cliente no hub
	hub.register <- client

	// Iniciar goroutines para ler e escrever no WebSocket
	go client.readPump()
	go client.writePump()

	log.Printf("Cliente %s (ID: %s) conectado", client.email, client.ID)
}
