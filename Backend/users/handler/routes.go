package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"users/models"

	"github.com/alexedwards/argon2id"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type Application struct {
	Models models.Models
}

func (app *Application) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	trimmedName := strings.TrimSpace(input.Name)
	fmt.Println("trimmed name: ", trimmedName)
	if trimmedName == "" {
		http.Error(w, "Name is required.", http.StatusBadRequest)
		return
	}

	nameParts := strings.Fields(trimmedName)
	if len(nameParts) < 2 {
		http.Error(w, "Please provide both first and last name", http.StatusBadRequest)
		return
	}

	if err := models.ValidateEduEmail(input.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := models.ValidatePassword(input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err := app.Models.UserModel.Insert(input.Id, input.Name, input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid Email or Password.", http.StatusBadRequest)
	}

}

func (app *Application) PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = app.Models.UserModel.Update(input.Email, input.Password)
	if err != nil {
		fmt.Print(err)
	}

}

func (app *Application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err = app.Models.UserModel.Delete(input.Email)
	if err != nil {
		fmt.Print(err)
	}
}

func (app *Application) DisplayUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func (app *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		http.Error(w, "Error occurred during password verification", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userId, err := app.Models.UserModel.GetUserIdByEmail(input.Email)
	if err != nil {
		http.Error(w, "Error occurred while fetching user ID", http.StatusInternalServerError)
		return
	}
	response_data := map[string]interface{}{
		"userId": userId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response_data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/signup", app.SignUpHandler)
	router.HandlerFunc(http.MethodPost, "/updatePassword", app.PasswordResetHandler)
	router.HandlerFunc(http.MethodPost, "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc(http.MethodPost, "/displayUser", app.DisplayUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.LoginHandler)

	// Apply CORS middleware
	return SetupCORS(router)
}

func SetupCORS(router http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true, // Log CORS-related issues, can be turned off for production
	})

	return c.Handler(router)
}
