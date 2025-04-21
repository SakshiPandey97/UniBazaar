// handler/handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"users/handler"
	"users/models"

	"github.com/alexedwards/argon2id"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHandlers(t *testing.T) {
	orig := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", orig)
	_ = os.Setenv("JWT_SECRET", "testsecret")

	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&models.User{})

	userModel := models.UserModel{DB: db}
	app := handler.Application{Models: models.Models{UserModel: userModel}}

	t.Run("SignUpHandler", func(t *testing.T) {
		body := map[string]interface{}{
			"id": 101, "name": "Ellie Williams", "email": "ellie@ufl.edu",
			"password": "SomeStrongPassword@123", "phone": "+15551234567",
		}
		req, _ := http.NewRequest(http.MethodPost, "/signup", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var created models.User
		assert.NoError(t, db.Where("email = ?", "ellie@ufl.edu").First(&created).Error)
		assert.False(t, created.Verified)
	})

	t.Run("VerifyEmailHandler", func(t *testing.T) {
		db.Create(&models.User{UserID: 202, Name: "Joel Miller", Email: "joel@ufl.edu", OTPCode: "123456"})
		body := map[string]string{"email": "joel@ufl.edu", "code": "123456"}
		req, _ := http.NewRequest(http.MethodPost, "/verifyEmail", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var u models.User
		_ = db.Where("email = ?", "joel@ufl.edu").First(&u)
		assert.True(t, u.Verified)
	})

	t.Run("ForgotPasswordHandler", func(t *testing.T) {
		db.Create(&models.User{UserID: 303, Name: "Tommy", Email: "tommy@ufl.edu", Verified: true})
		req, _ := http.NewRequest(http.MethodGet, "/forgotPassword?email=tommy@ufl.edu", nil)
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.User
		_ = db.Where("email = ?", "tommy@ufl.edu").First(&updated)
		assert.NotEmpty(t, updated.OTPCode)
	})

	t.Run("UpdatePasswordHandler", func(t *testing.T) {
		db.Create(&models.User{UserID: 404, Name: "Sarah", Email: "sarah@ufl.edu", OTPCode: "654321", Verified: true})
		body := map[string]string{
			"email": "sarah@ufl.edu", "otp_code": "654321", "new_password": "BrandNewPassword@99",
		}
		req, _ := http.NewRequest(http.MethodPost, "/updatePassword", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.User
		_ = db.Where("email = ?", "sarah@ufl.edu").First(&updated)
		assert.Empty(t, updated.OTPCode)
	})

	t.Run("DisplayUserHandler", func(t *testing.T) {
		db.Create(&models.User{UserID: 505, Name: "Henry", Email: "henry@ufl.edu"})
		req, _ := http.NewRequest(http.MethodGet, "/displayUser/505", nil)
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp models.User
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "Henry", resp.Name)
	})

	t.Run("UpdateNameHandler", func(t *testing.T) {
		raw := "Mychemicalromance!567"
		params := &argon2id.Params{
			Memory:      65536,
			Iterations:  2,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		}
		hash, _ := argon2id.CreateHash(raw, params)
		db.Create(&models.User{UserID: 606, Name: "Bill", Email: "bill@ufl.edu", Password: hash, Verified: true})

		body := map[string]string{"email": "bill@ufl.edu", "password": raw, "newName": "Billy"}
		req, _ := http.NewRequest(http.MethodPost, "/updateName", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.User
		_ = db.Where("email = ?", "bill@ufl.edu").First(&updated)
		assert.Equal(t, "Billy", updated.Name)
	})

	t.Run("UpdatePhoneHandler", func(t *testing.T) {
		raw := "SecretPhone99!"
		params := &argon2id.Params{
			Memory:      65536,
			Iterations:  2,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		}
		hash, _ := argon2id.CreateHash(raw, params)
		db.Create(&models.User{UserID: 707, Name: "Maria", Email: "maria@ufl.edu", Password: hash, Phone: "5551234567", Verified: true})

		body := map[string]string{"email": "maria@ufl.edu", "password": raw, "newPhone": "5559998888"}
		req, _ := http.NewRequest(http.MethodPost, "/updatePhone", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.User
		_ = db.Where("email = ?", "maria@ufl.edu").First(&updated)
		assert.Equal(t, "5559998888", updated.Phone)
	})

	t.Run("DeleteUserHandler", func(t *testing.T) {
		db.Create(&models.User{UserID: 808, Name: "David", Email: "david@ufl.edu"})
		body := map[string]string{"email": "david@ufl.edu"}
		req, _ := http.NewRequest(http.MethodPost, "/deleteUser", toJSON(body))
		rec := httptest.NewRecorder()
		app.Routes().ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var gone models.User
		assert.Error(t, db.Where("email = ?", "david@ufl.edu").First(&gone).Error)
	})

	t.Log("All handler tests completed:", time.Now())
}

func toJSON(v interface{}) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
