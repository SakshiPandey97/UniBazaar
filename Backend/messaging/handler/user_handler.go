package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"messaging/repository"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type UserHandler struct {
	UserRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding users:", err)
		http.Error(w, "Error encoding users", http.StatusInternalServerError)
	}
}

type SyncUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) SyncUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload SyncUserRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&reqPayload)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON or is incomplete"
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		case errors.As(err, &invalidUnmarshalError):
			log.Printf("Internal Server Error: Invalid argument passed to json.Unmarshal: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	if reqPayload.ID == 0 || reqPayload.Name == "" || reqPayload.Email == "" {
		http.Error(w, "Missing required fields: id, name, email", http.StatusBadRequest)
		return
	}

	exists, err := h.UserRepo.UserExists(reqPayload.ID)
	if err != nil {
		log.Printf("Error checking if user %d exists during sync: %v", reqPayload.ID, err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Printf("User %d already exists in messaging DB. Sync successful.", reqPayload.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User already synchronized"})
		return
	}

	log.Printf("User %d does not exist in messaging DB. Attempting to add...", reqPayload.ID)
	err = h.UserRepo.AddUser(reqPayload.ID, reqPayload.Name, reqPayload.Email)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			log.Printf("User %d was added concurrently. Sync successful.", reqPayload.ID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "User already synchronized (concurrent add)"})
		} else {
			// Other database error
			log.Printf("Error adding user %d during sync: %v", reqPayload.ID, err)
			http.Error(w, "Failed to add user during sync", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully added user %d during sync.", reqPayload.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User synchronized successfully"})
}
