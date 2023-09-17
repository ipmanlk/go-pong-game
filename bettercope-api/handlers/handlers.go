package handlers

import (
	"encoding/json"
	"ipmanlk/bettercopelk/models"
	"ipmanlk/bettercopelk/search"
	"log"
	"net/http"
	"strings"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	// Get the query from the URL query string
	query := strings.TrimSpace(r.URL.Query().Get("query"))

	// If the query is empty after triming, return a bad request status
	if query == "" {
		http.Error(w, "Missing query", http.StatusBadRequest)
		return
	}

	// Set response headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a channel to send scraped results to the SSE client
	resultsChan := make(chan []models.SearchResult, 3)

	search.SearchSites(query, resultsChan)

	for results := range resultsChan {
		// Serialize the result as JSON
		resultJSON, err := json.Marshal(results)
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			continue
		}

		// Wrap the JSON result in an SSE message
		sseMessage := "event: results\ndata: " + string(resultJSON) + "\n\n"

		_, err = w.Write([]byte(sseMessage))
		if err != nil {
			log.Println("Error writing SSE message:", err)
			return
		}

		w.(http.Flusher).Flush() // Flush the response to send it immediately
	}

	// Signal the end of the SSE stream
	endMessage := "event: end\ndata: end\n\n"
	_, err := w.Write([]byte(endMessage))
	if err != nil {
		log.Println("Error writing SSE end message:", err)
	}
}
