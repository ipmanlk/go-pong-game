package main

import (
	"net/http"

	"ipmanlk/gochat/chatroom"
	"ipmanlk/gochat/database"
	"ipmanlk/gochat/handlers"
)

func main() {
	database.InitDatabase()

	chatRoom := chatroom.NewChatRoom()
	go chatRoom.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocketConnection(chatRoom, w, r)
	})
	http.HandleFunc("/messages", handlers.GetLastMessagesHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Failed to start the server: " + err.Error())
	}
}
