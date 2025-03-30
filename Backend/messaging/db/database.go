package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	// connStr := "user=postgres password=postgres2025 dbname=messaging sslmode=disable"
	connStr := "user=tanmay password=postgres2025 dbname=unibazaar-messaging-db host=unibazaar-messaging-db.ca7004ecmjyv.us-east-1.rds.amazonaws.com port=5432 sslmode=require"
	//
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
