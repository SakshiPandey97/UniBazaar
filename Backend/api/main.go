package main

import (
	"Backend/data/models"
	"fmt"
	"net/http"
)

type application struct {
	Models models.Models
}

func main() {
	dsn := "postgres://postgres:admin@localhost/unibazaar?sslmode=disable"
	app := application{}
	conn := Connect(dsn)
	app.Models = models.NewModels(conn)
	fmt.Println("connected to database")

	srv := http.Server{
		Addr:    ":4000",
		Handler: app.routes(),
	}
	fmt.Printf("app running on port %d", 4000)
	_ = srv.ListenAndServe()
}
