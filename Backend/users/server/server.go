package server

import (
	"fmt"
	"net/http"
	config "users/config"
	handler "users/handler"
	models "users/models"
)

func InitServer() {
	dsn := "postgres://postgres:Shubham16!@localhost/unibazar?sslmode=disable"
	app := handler.Application{}
	conn := config.Connect(dsn)
	app.Models = models.NewModels(conn)
	fmt.Println("connected to database")

	srv := http.Server{
		Addr:    ":4000",
		Handler: app.Routes(),
	}
	fmt.Printf("app running on port %d", 4000)
	_ = srv.ListenAndServe()
}
