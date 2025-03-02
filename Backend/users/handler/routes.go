package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"users/models"

	"github.com/alexedwards/argon2id"
	"github.com/julienschmidt/httprouter"
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
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Sign-up successful. Please check your email for the OTP.")
}

func (app *Application) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email   string `json:"email"`
		OTPCode string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user.OTPCode != input.OTPCode {
		user.FailedResetAttempts++
		if user.FailedResetAttempts >= 3 {
			_ = app.Models.UserModel.SendSecurityAlert(user)
			user.FailedResetAttempts = 0
			user.OTPCode = ""
		}
		_ = app.Models.UserModel.SaveUser(user)
		http.Error(w, "Invalid OTP code.", http.StatusUnauthorized)
		return
	}
	user.Verified = true
	user.FailedResetAttempts = 0
	user.OTPCode = ""
	if err := app.Models.UserModel.UpdateVerificationStatus(user); err != nil {
		http.Error(w, "Failed to update verification status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Email verified successfully!")
}

func (app *Application) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := app.Models.UserModel.InitiatePasswordReset(input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Reset code sent. Check your email.")
}

func (app *Application) VerifyResetCodeHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email       string `json:"email"`
		OTPCode     string `json:"otp_code"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := app.Models.UserModel.VerifyResetCodeAndSetNewPassword(input.Email, input.OTPCode, input.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password reset successful.")
}

func (app *Application) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email       string `json:"email"`
		OTPCode     string `json:"otp_code"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	app.Models.UserModel.VerifyResetCodeAndSetNewPassword(input.Email, input.OTPCode, input.NewPassword)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password updated successfully.")
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
	router.HandlerFunc("POST", "/signup", app.SignUpHandler)
	router.HandlerFunc("POST", "/verifyEmail", app.VerifyEmailHandler)
	router.HandlerFunc("POST", "/forgotPassword", app.ForgotPasswordHandler)
	router.HandlerFunc("POST", "/verifyResetCode", app.VerifyResetCodeHandler)
	router.HandlerFunc("POST", "/updatePassword", app.UpdatePasswordHandler)
	router.HandlerFunc("POST", "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc("POST", "/displayUser", app.DisplayUserHandler)
	router.HandlerFunc("POST", "/login", app.LoginHandler)
	return router
}

//add better errors like for the json stuff, or the email being sent.
//fix sign up if no otp is sent :/
