package handlers

import (
	"encoding/json"
	"ipmanlk/gochat/database"
	"net/http"
)

func GetLastMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages := database.GetLastNMessages(100)
	response, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
