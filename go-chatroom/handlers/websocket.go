package handlers

import (
	"ipmanlk/gochat/chatroom"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocketConnection(cr *chatroom.ChatRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during WebSocket upgrade:", err)
		return
	}
	client := &chatroom.Client{Conn: conn, Send: make(chan []byte, 256)}
	cr.Register <- client

	go client.Write()
	go client.Read(cr)
}
