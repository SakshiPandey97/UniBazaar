package handler

import (
	"encoding/json"
	"log"
	"messaging/repository"
	"net/http"
)

// UserHandler struct contains a reference to the user repository
type UserHandler struct {
	UserRepo *repository.UserRepository
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// GetUsersHandler handles the request to fetch all users
func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	// Set the response header and return the user data as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding users:", err)
		http.Error(w, "Error encoding users", http.StatusInternalServerError)
	}
}
