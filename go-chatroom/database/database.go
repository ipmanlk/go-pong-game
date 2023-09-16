package database

import (
	"database/sql"
	"ipmanlk/gochat/common"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "chat.db")
	if err != nil {
		log.Fatal(err)
	}
	createMessageTable()
}

func createMessageTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS messages (username TEXT, content TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}

func SaveMessage(username, content string) {
	_, err := db.Exec("INSERT INTO messages (username, content) VALUES (?, ?)", username, content)
	if err != nil {
		log.Println("Failed to insert message:", err)
	}
}

func GetLastNMessages(n int) []common.Message {
	rows, err := db.Query("SELECT username, content FROM messages ORDER BY rowid DESC LIMIT ?", n)
	if err != nil {
		log.Println("Failed to fetch messages:", err)
		return nil
	}
	defer rows.Close()

	var messages []common.Message
	for rows.Next() {
		var msg common.Message
		if err := rows.Scan(&msg.Username, &msg.Content); err != nil {
			log.Println("Failed to scan message:", err)
			return messages
		}
		messages = append(messages, msg)
	}
	return messages
}
