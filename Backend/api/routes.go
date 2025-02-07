package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to UniBazaar!"))
}

func (app *application) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dec := json.NewDecoder(r.Body)
	dec.Decode(&input)

	err := app.Models.UserModel.Insert(input.Id, input.Name, input.Email, input.Password)
	if err != nil {
		fmt.Print(err)
	}

}
func (app *application) PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) DisplayUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/signup", app.SignUpHandler)
	router.HandlerFunc(http.MethodPost, "/updatePassword", app.PasswordResetHandler)
	router.HandlerFunc(http.MethodPost, "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc(http.MethodPost, "/displayUser", app.DisplayUserHandler)
	return router
}

func testing() {
	fmt.Print("hello")
}
