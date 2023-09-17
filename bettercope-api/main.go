package main

import (
	"ipmanlk/bettercopelk/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/search", handlers.HandleSearch)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
