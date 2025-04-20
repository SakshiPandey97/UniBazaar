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

// --- Add SyncUser Handler ---

// SyncUserRequest defines the expected JSON payload for syncing
type SyncUserRequest struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SyncUserHandler handles the POST /api/users/sync request
// It checks if the user exists before attempting to add them.
func (h *UserHandler) SyncUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload SyncUserRequest
	// Limit request body size for security
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576) // 1MB limit

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent unexpected fields in JSON

	err := decoder.Decode(&reqPayload)
	if err != nil {
		// --- Detailed JSON error handling ---
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
		case err.Error() == "http: request body too large": // Check specific error string for MaxBytesReader
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		case errors.As(err, &invalidUnmarshalError):
			// This indicates a problem with the Go struct passed to Decode
			log.Printf("Internal Server Error: Invalid argument passed to json.Unmarshal: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
		// --- End of JSON error handling ---
	}

	// Basic validation on required fields
	if reqPayload.ID == 0 || reqPayload.Name == "" || reqPayload.Email == "" {
		http.Error(w, "Missing required fields: id, name, email", http.StatusBadRequest)
		return
	}

	// Check if user already exists in the messaging database
	exists, err := h.UserRepo.UserExists(reqPayload.ID)
	if err != nil {
		log.Printf("Error checking if user %d exists during sync: %v", reqPayload.ID, err)
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	// If user already exists, return success (idempotent operation)
	if exists {
		log.Printf("User %d already exists in messaging DB. Sync successful.", reqPayload.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User already synchronized"})
		return
	}

	// User does not exist, attempt to add them
	log.Printf("User %d does not exist in messaging DB. Attempting to add...", reqPayload.ID)
	// Assuming AddUser function exists in your repository
	err = h.UserRepo.AddUser(reqPayload.ID, reqPayload.Name, reqPayload.Email)
	if err != nil {
		// Check if the error is a unique constraint violation (e.g., added concurrently)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // PostgreSQL specific code for unique_violation
			log.Printf("User %d was added concurrently. Sync successful.", reqPayload.ID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // Treat as success if it exists now
			json.NewEncoder(w).Encode(map[string]string{"message": "User already synchronized (concurrent add)"})
		} else {
			// Other database error
			log.Printf("Error adding user %d during sync: %v", reqPayload.ID, err)
			http.Error(w, "Failed to add user during sync", http.StatusInternalServerError)
		}
		return
	}

	// Successfully added the user
	log.Printf("Successfully added user %d during sync.", reqPayload.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Indicate resource was created
	json.NewEncoder(w).Encode(map[string]string{"message": "User synchronized successfully"})
}
