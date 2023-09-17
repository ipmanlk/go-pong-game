package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/bettercopelk/models"
	"ipmanlk/bettercopelk/search"
	"os"
)

func main() {
	query := "life"
	resultsChan := make(chan []models.SearchResult)

	search.SearchSites(query, resultsChan)

	results := []models.SearchResult{}

	for searchResults := range resultsChan {
		fmt.Println(searchResults)
		results = append(results, searchResults...)
	}

	// Wrilte results to a file using the results slice
	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(results); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Done!")
}
