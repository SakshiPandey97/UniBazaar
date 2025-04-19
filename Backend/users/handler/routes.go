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

// Global in-memory map to store revoked token IDs (jti).
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
		fmt.Printf("SignUpHandler error: could not create user: %v\n", err)
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
		fmt.Printf("VerifyEmailHandler error: user not found: %v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if user.OTPCode != input.Code {
		user.FailedResetAttempts++
		if user.FailedResetAttempts >= 5 {
			fmt.Printf("VerifyEmailHandler: user %s reached 5 OTP failures, deleting account\n", user.Email)
			_ = app.Models.UserModel.Delete(user.Email)
			http.Error(w, "Too many OTP attempts. Your account has been removed. Please sign up again.", http.StatusUnauthorized)
			return
		}
		if user.FailedResetAttempts >= 3 {
			_ = app.Models.UserModel.SendSecurityAlert(user)
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
		fmt.Printf("VerifyEmailHandler error: %v\n", err)
		http.Error(w, "failed to update verification status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Email verified successfully!")
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
		fmt.Printf("ResendOTPHandler error: %v\n", err)
		http.Error(w, fmt.Sprintf("failed to resend OTP: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OTP resent successfully. Check your e‑mail.")
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
		fmt.Printf("ForgotPasswordHandler error: %v\n", err)
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
		fmt.Printf("UpdatePasswordHandler error: %v\n", err)
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
		fmt.Printf("DeleteUserHandler error: %v\n", err)
		http.Error(w, fmt.Sprintf("delete failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted.")
}

func (app *Application) DisplayUserHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	idStr := params.ByName("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := app.Models.UserModel.ReadByID(userID)
	if err != nil {
		fmt.Printf("DisplayUserHandler error: %v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
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
		fmt.Printf("LoginHandler error: %v\n", err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	if !user.Verified {
		fmt.Printf("LoginHandler: user %s not verified\n", user.Email)
		http.Error(w, "account not verified, please verify your email first", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		fmt.Printf("LoginHandler error verifying password: %v\n", err)
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := utils.GenerateJWT(*user)
	if err != nil {
		fmt.Printf("LoginHandler error: %v\n", err)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	userID, err := app.Models.UserModel.GetUserIdByEmail(input.Email)
	if err != nil {
		fmt.Printf("LoginHandler error: %v\n", err)
		http.Error(w, "failed to fetch user ID", http.StatusInternalServerError)
		return
	}
	responseData := map[string]interface{}{
		"userId": userID,
		"token":  tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
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
		fmt.Printf("UpdateNameHandler error: %v\n", err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		fmt.Printf("UpdateNameHandler verify pwd: %v\n", err)
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := app.Models.UserModel.UpdateName(input.Email, input.NewName); err != nil {
		fmt.Printf("UpdateNameHandler error: %v\n", err)
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
		fmt.Printf("UpdatePhoneHandler error: %v\n", err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		fmt.Printf("UpdatePhoneHandler verify pwd: %v\n", err)
		http.Error(w, "error verifying password", http.StatusInternalServerError)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := app.Models.UserModel.UpdatePhone(input.Email, input.NewPhone); err != nil {
		fmt.Printf("UpdatePhoneHandler error: %v\n", err)
		http.Error(w, fmt.Sprintf("update phone failed: %v", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Phone updated successfully.")
}

// LogoutHandler revokes the current token
func (app *Application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	parts := strings.Split(bearer, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "invalid authorization header", http.StatusBadRequest)
		return
	}
	authToken := parts[1]

	jwtToken, err := utils.ParseJWT(authToken)
	if err != nil || !jwtToken.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "invalid token claims", http.StatusUnauthorized)
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		http.Error(w, "token missing jti claim", http.StatusUnauthorized)
		return
	}

	revokedTokens[jti] = true
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Logout successful, token revoked.")
}

func (app *Application) VerifyJWTHandler(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	parts := strings.Split(bearer, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		fmt.Println("invalid authorization header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authToken := parts[1]
	fmt.Printf("the authorization token passed by the user is: %s \n", authToken)

	jwtToken, err := utils.ParseJWT(authToken)
	if err != nil || !jwtToken.Valid {
		fmt.Println("error parsing JWT or token invalid:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if jti is revoked
	jti, ok := claims["jti"].(string)
	if ok && revokedTokens[jti] {
		fmt.Println("token has been revoked")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userMap, ok := claims["user"].(map[string]interface{})
	if !ok {
		fmt.Println("error retrieving user claims")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userClaim := models.User{
		Name:  fmt.Sprintf("%v", userMap["Name"]),
		Email: fmt.Sprintf("%v", userMap["Email"]),
		Phone: fmt.Sprintf("%v", userMap["Phone"]),
	}

	fmt.Printf("JWT is valid. User from token: %#v\n", userClaim)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf("Token valid. User: %v", userClaim)))
}

func (app *Application) GetJWTHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user := models.CreateUser(input.Name, input.Email, input.Phone)
	authToken, err := utils.GenerateJWT(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+authToken)
	_, _ = w.Write([]byte("JWT generated successfully!"))
}

func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	/*  AUTH  */
	router.HandlerFunc("POST", "/signup", app.SignUpHandler)
	router.HandlerFunc("POST", "/verifyEmail", app.VerifyEmailHandler)
	router.HandlerFunc("POST", "/resendOtp", app.ResendOTPHandler)

	/*  PASSWORD RESET  */
	router.HandlerFunc("POST", "/forgotPassword", app.ForgotPasswordHandler)
	router.HandlerFunc("POST", "/updatePassword", app.UpdatePasswordHandler)

	/*  USER CRUD  */
	router.HandlerFunc("POST", "/deleteUser", app.DeleteUserHandler)
	router.HandlerFunc("GET", "/displayUser/:id", app.DisplayUserHandler)
	router.HandlerFunc("POST", "/updateName", app.UpdateNameHandler)
	router.HandlerFunc("POST", "/updatePhone", app.UpdatePhoneHandler)

	/*  TOKEN  */
	router.HandlerFunc("POST", "/login", app.LoginHandler)
	router.HandlerFunc("POST", "/logout", app.LogoutHandler)
	router.HandlerFunc("POST", "/getjwt", app.GetJWTHandler)
	router.HandlerFunc("GET", "/verifyjwt", app.VerifyJWTHandler)

	return app.enableCORS(router)
}
