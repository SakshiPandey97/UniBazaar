package routes

import (
	"messaging/handler"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, messageHandler *handler.MessageHandler) {
	router.HandleFunc("/ws", messageHandler.HandleWebSocket)
}
