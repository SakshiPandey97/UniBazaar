package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	config "users/config"
	handler "users/handler"
	models "users/models"

	"github.com/joho/godotenv" // go get github.com/joho/godotenv
)

func InitServer() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found relying on real environment variables")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN environment variable not set")
	}

	conn := config.Connect(dsn)

	app := handler.Application{
		Models: models.NewModels(conn),
	}

	fmt.Println("connected to database")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	//	srv := http.Server{
	//		Addr:    ":8080",
	//		Handler: app.Routes(),
	//}

	fmt.Println("app running on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
