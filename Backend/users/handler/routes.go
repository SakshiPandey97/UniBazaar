package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"users/models"
	"users/utils"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

var revokedTokens = make(map[string]bool)

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
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	nameParts := strings.Fields(strings.TrimSpace(input.Name))
	if len(nameParts) < 2 {
		http.Error(w, "please provide both first and last name", http.StatusBadRequest)
		return
	}
	if err := app.Models.UserModel.Insert(input.Id, input.Name, input.Email, input.Password, input.Phone); err != nil {
		http.Error(w, fmt.Sprintf("could not create user: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Sign‑up successful. Check your email for OTP.")
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
		switch {
		case user.FailedResetAttempts >= 5:
			_ = app.Models.UserModel.Delete(user.Email)
		case user.FailedResetAttempts >= 3:
			_ = app.Models.UserModel.SendSecurityAlert(user)
			user.OTPCode = ""
			_ = app.Models.UserModel.SaveUser(user)
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]bool{"verified": false})
		return
	}
	user.Verified, user.FailedResetAttempts, user.OTPCode = true, 0, ""
	if err := app.Models.UserModel.UpdateVerificationStatus(user); err != nil {
		http.Error(w, "failed to update verification status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"verified": true})
}

func (app *Application) ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}
	if err := app.Models.UserModel.ResendOTP(input.Email); err != nil {
		http.Error(w, fmt.Sprintf("failed to resend OTP: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OTP resent successfully. Check your e‑mail.")
}

func (app *Application) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if strings.TrimSpace(email) == "" {
		http.Error(w, "email query param required", http.StatusBadRequest)
		return
	}
	if err := app.Models.UserModel.InitiatePasswordReset(email); err != nil {
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
	if err := app.Models.UserModel.VerifyResetCodeAndSetNewPassword(input.Email, input.OTPCode, input.NewPassword); err != nil {
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
	if err := app.Models.UserModel.Delete(input.Email); err != nil {
		http.Error(w, fmt.Sprintf("delete failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted.")
}

func (app *Application) DisplayUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := app.Models.UserModel.ReadByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
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
	if !user.Verified {
		http.Error(w, "account not verified, please verify your email first", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil || !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := utils.GenerateJWT(*user)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	userID, _ := app.Models.UserModel.GetUserIdByEmail(input.Email)
	resp := map[string]interface{}{
		"userId": userID,
		"token":  tokenString,
		"name":   user.Name,
		"email":  user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
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
	match, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
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
	match, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password)
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

func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	bearer := strings.Split(r.Header.Get("Authorization"), " ")
	if len(bearer) != 2 || bearer[0] != "Bearer" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}
	jwtToken, err := utils.ParseJWT(bearer[1])
	if err != nil || !jwtToken.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if jti, ok := claims["jti"].(string); ok {
		revokedTokens[jti] = true
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Logout successful, token revoked.")
}

func (app *Application) VerifyJWTHandler(w http.ResponseWriter, r *http.Request) {
	bearer := strings.Split(r.Header.Get("Authorization"), " ")
	if len(bearer) != 2 || bearer[0] != "Bearer" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jwtToken, err := utils.ParseJWT(bearer[1])
	if err != nil || !jwtToken.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if jti, ok := claims["jti"].(string); ok && revokedTokens[jti] {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userMap := claims["user"].(map[string]interface{})
	userClaim := models.User{
		Name:  fmt.Sprintf("%v", userMap["Name"]),
		Email: fmt.Sprintf("%v", userMap["Email"]),
		Phone: fmt.Sprintf("%v", userMap["Phone"]),
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("Token valid. User: %v", userClaim)))
}

func (app *Application) GetJWTHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user := models.CreateUser(input.Name, input.Email, input.Phone)
	token, _ := utils.GenerateJWT(*user)
	w.Header().Set("Authorization", "Bearer "+token)
	_, _ = w.Write([]byte("JWT generated successfully!"))
}

func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) Routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("OK")) })

	router.HandlerFunc("POST", "/signup", app.SignUpHandler)
	router.HandlerFunc("POST", "/verifyEmail", app.VerifyEmailHandler)
	router.HandlerFunc("POST", "/resendOtp", app.ResendOTPHandler)

	router.HandlerFunc("GET", "/forgotPassword", app.ForgotPasswordHandler)
	router.HandlerFunc("POST", "/updatePassword", app.UpdatePasswordHandler)

	router.HandlerFunc("POST", "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc("GET", "/displayUser/:id", app.DisplayUserHandler)
	router.HandlerFunc("POST", "/updateName", app.UpdateNameHandler)
	router.HandlerFunc("POST", "/updatePhone", app.UpdatePhoneHandler)

	router.HandlerFunc("POST", "/login", app.LoginHandler)
	router.HandlerFunc("POST", "/logout", app.LogoutHandler)
	router.HandlerFunc("POST", "/getjwt", app.GetJWTHandler)
	router.HandlerFunc("GET", "/verifyjwt", app.VerifyJWTHandler)

	return app.enableCORS(router)
}
