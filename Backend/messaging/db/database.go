package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}
	connStr := os.Getenv("CHAT_DB_URI")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Database not responding:", err)
	}
	fmt.Println("Connected to database")
	return db
}
