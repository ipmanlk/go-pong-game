package chatroom

import (
	"encoding/json"
	"ipmanlk/gochat/common"
	"ipmanlk/gochat/database"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type ChatRoom struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func (c *Client) Write() {
	defer c.Conn.Close()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}

func (c *Client) Read(cr *ChatRoom) {
	defer func() {
		cr.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgBytes, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var message common.Message
		if err = json.Unmarshal(msgBytes, &message); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		formattedMessage, err := json.Marshal(message)
		if err != nil {
			log.Println("Error encoding JSON:", err)
			continue
		}

		cr.broadcast <- formattedMessage
		database.SaveMessage(message.Username, message.Content)
	}
}

func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.Register:
			cr.mu.Lock()
			cr.clients[client] = true
			cr.mu.Unlock()
		case client := <-cr.unregister:
			cr.mu.Lock()
			if _, ok := cr.clients[client]; ok {
				close(client.Send)
				delete(cr.clients, client)
			}
			cr.mu.Unlock()
		case message := <-cr.broadcast:
			cr.mu.Lock()
			for client := range cr.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(cr.clients, client)
				}
			}
			cr.mu.Unlock()
		}
	}
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
