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
		Phone    string `json:"phone"`
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	trimmedName := strings.TrimSpace(input.Name)
	nameParts := strings.Fields(trimmedName)
	if len(nameParts) < 2 {
		http.Error(w, "please provide both first and last name", http.StatusBadRequest)
		return
	}

	err := app.Models.UserModel.Insert(input.Id, input.Name, input.Email, input.Password, input.Phone)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create user: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Sign-up successful. Check your email for OTP.")
}


func (app *Application) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if user.OTPCode != input.Code {
		user.FailedResetAttempts++
		if user.FailedResetAttempts >= 3 {
			_ = app.Models.UserModel.SendSecurityAlert(user)
			user.FailedResetAttempts = 0
			user.OTPCode = ""
		}
		_ = app.Models.UserModel.SaveUser(user)
		http.Error(w, "invalid OTP code", http.StatusUnauthorized)
		return
	}
	user.Verified = true
	user.FailedResetAttempts = 0
	user.OTPCode = ""
	if err := app.Models.UserModel.UpdateVerificationStatus(user); err != nil {
		http.Error(w, "failed to update verification status", http.StatusInternalServerError)
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
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	err := app.Models.UserModel.InitiatePasswordReset(input.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to initiate reset: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Reset code sent. Check your email.")
}



func (app *Application) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email       string `json:"email"`
		OTPCode     string `json:"otp_code"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	err := app.Models.UserModel.VerifyResetCodeAndSetNewPassword(input.Email, input.OTPCode, input.NewPassword)
	if err != nil {
		http.Error(w, fmt.Sprintf("update password failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Password updated successfully.")
}

func (app *Application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	err := app.Models.UserModel.Delete(input.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("delete failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted.")
}

func (app *Application) DisplayUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
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
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	userID, err := app.Models.UserModel.GetUserIdByEmail(input.Email)
	if err != nil {
		http.Error(w, "failed to fetch user ID", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	responseData := map[string]interface{}{"userId": userID}
	json.NewEncoder(w).Encode(responseData)
	w.WriteHeader(http.StatusOK)
}

func (app *Application) UpdateNameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		NewName  string `json:"newName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := app.Models.UserModel.UpdateName(input.Email, input.NewName); err != nil {
		http.Error(w, fmt.Sprintf("update name failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Name updated successfully.")
}

func (app *Application) UpdatePhoneHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		NewPhone string `json:"newPhone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.Read(input.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := app.Models.UserModel.UpdatePhone(input.Email, input.NewPhone); err != nil {
		http.Error(w, fmt.Sprintf("update phone failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Phone updated successfully.")
}

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc("POST", "/signup", app.SignUpHandler)
	router.HandlerFunc("POST", "/verifyEmail", app.VerifyEmailHandler)
	router.HandlerFunc("POST", "/forgotPassword", app.ForgotPasswordHandler)
	router.HandlerFunc("POST", "/updatePassword", app.UpdatePasswordHandler)
	router.HandlerFunc("POST", "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc("POST", "/displayUser", app.DisplayUserHandler)
	router.HandlerFunc("POST", "/login", app.LoginHandler)
	router.HandlerFunc("POST", "/updateName", app.UpdateNameHandler)
	router.HandlerFunc("POST", "/updatePhone", app.UpdatePhoneHandler)

	return router
}
