package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client representa um usuário conectado ao WebSocket
type Client struct {
    ID     string
    email string 
    Conn   *websocket.Conn
    Pool   *Hub
    Send   chan Message
}


// Hub gerencia todas as conexões WebSocket
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// Message representa a estrutura da mensagem enviada pelo WebSocket
type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Type    string `json:"type"` // "text", "image", etc.
}

// NewHub cria uma nova instância do Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Start inicializa o loop do Hub para gerenciar conexões e mensagens
func (h *Hub) Start() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Cliente %s conectado", client.ID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				client.Conn.Close()
				log.Printf("Cliente %s desconectado", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// readPump lê mensagens do cliente e as envia para o hub
func (c *Client) readPump() {
	defer func() {
		c.Pool.unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}
		msg.Sender = c.email
		c.Pool.broadcast <- msg
	}
}

// writePump envia mensagens do hub para o cliente
func (c *Client) writePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			log.Printf("Erro ao escrever mensagem: %v", err)
			break
		}
	}
}