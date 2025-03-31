package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "testsecret")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	userModel := models.UserModel{DB: db}
	app := handler.Application{
		Models: models.Models{UserModel: userModel},
	}

	t.Run("SignUpHandler", func(t *testing.T) {
		body := map[string]interface{}{
			"id":       101,
			"name":     "Ellie Williams",
			"email":    "ellie@ufl.edu",
			"password": "SomeStrongPassword@123",
			"phone":    "+15551234567",
		}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var createdUser models.User
		err = db.Where("email = ?", "ellie@ufl.edu").First(&createdUser).Error
		assert.NoError(t, err)
		assert.False(t, createdUser.Verified)
	})

	t.Run("VerifyEmailHandler", func(t *testing.T) {
		user := models.User{
			UserID:   202,
			Name:     "Joel Miller",
			Email:    "joel@ufl.edu",
			Phone:    "5559876543",
			Verified: false,
			OTPCode:  "123456",
		}
		_ = db.Create(&user)

		body := map[string]string{
			"email": "joel@ufl.edu",
			"code":  "123456",
		}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/verifyEmail", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var verifiedUser models.User
		_ = db.Where("email = ?", "joel@ufl.edu").First(&verifiedUser)
		assert.True(t, verifiedUser.Verified)
	})

	t.Run("ForgotPasswordHandler", func(t *testing.T) {
		user := models.User{
			UserID:   303,
			Name:     "Tommy Miller",
			Email:    "tommy@ufl.edu",
			Password: "dummy-hash-here",
			Verified: true,
		}
		_ = db.Create(&user)

		body := map[string]string{"email": "tommy@ufl.edu"}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/forgotPassword", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updated models.User
		_ = db.Where("email = ?", "tommy@ufl.edu").First(&updated)
		assert.NotEmpty(t, updated.OTPCode)
	})

	t.Run("UpdatePasswordHandler", func(t *testing.T) {
		user := models.User{
			UserID:   404,
			Name:     "Sarah",
			Email:    "sarah@ufl.edu",
			Password: "old-hash-here",
			Verified: true,
			OTPCode:  "654321",
		}
		_ = db.Create(&user)

		body := map[string]string{
			"email":        "sarah@ufl.edu",
			"otp_code":     "654321",
			"new_password": "BrandNewPassword@99",
		}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/updatePassword", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updated models.User
		_ = db.Where("email = ?", "sarah@ufl.edu").First(&updated)
		assert.NotEqual(t, "old-hash-here", updated.Password)
		assert.Empty(t, updated.OTPCode)
	})

	t.Run("DisplayUserHandler", func(t *testing.T) {
		user := models.User{
			UserID: 505,
			Name:   "Henry",
			Email:  "henry@ufl.edu",
			Phone:  "5557772222",
		}
		_ = db.Create(&user)

		body := map[string]string{"email": "henry@ufl.edu"}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/displayUser", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp models.User
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "Henry", resp.Name)
		assert.Equal(t, "henry@ufl.edu", resp.Email)
	})

	t.Run("UpdateNameHandler", func(t *testing.T) {
		rawPass := "MySecret42!"
		hashed, _ := argon2id.CreateHash(rawPass, &argon2id.Params{
			Memory:      64 * 1024,
			Iterations:  2,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		})

		user := models.User{
			UserID:   606,
			Name:     "Bill",
			Email:    "bill@ufl.edu",
			Password: hashed,
			Verified: true,
		}
		_ = db.Create(&user)

		body := map[string]string{
			"email":    "bill@ufl.edu",
			"password": rawPass,
			"newName":  "Bill the Survivor",
		}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/updateName", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updated models.User
		_ = db.Where("email = ?", "bill@ufl.edu").First(&updated)
		assert.Equal(t, "Bill the Survivor", updated.Name)
	})

	t.Run("UpdatePhoneHandler", func(t *testing.T) {
		rawPass := "SecretPhone99!"
		hashed, _ := argon2id.CreateHash(rawPass, &argon2id.Params{
			Memory:      64 * 1024,
			Iterations:  2,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   32,
		})

		user := models.User{
			UserID:   707,
			Name:     "Maria",
			Email:    "maria@ufl.edu",
			Password: hashed,
			Verified: true,
			Phone:    "5551234567",
		}
		_ = db.Create(&user)

		body := map[string]string{
			"email":    "maria@ufl.edu",
			"password": rawPass,
			"newPhone": "5559998888",
		}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/updatePhone", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var updated models.User
		_ = db.Where("email = ?", "maria@ufl.edu").First(&updated)
		assert.Equal(t, "5559998888", updated.Phone)
	})

	t.Run("DeleteUserHandler", func(t *testing.T) {
		user := models.User{
			UserID: 808,
			Name:   "David",
			Email:  "david@ufl.edu",
			Phone:  "5550001111",
		}
		_ = db.Create(&user)

		body := map[string]string{"email": "david@ufl.edu"}
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/deleteUser", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		app.Routes().ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var shouldNotExist models.User
		err := db.Where("email = ?", "david@ufl.edu").First(&shouldNotExist).Error
		assert.Error(t, err)
	})

	fmt.Println("All handler tests completed successfully:", time.Now())
}
