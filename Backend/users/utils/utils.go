package utils

import (
	"fmt"
	"os"
	"reflect"
	"time"
	"users/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// StructToMap maps a struct's fields into a map
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()
		result[fieldName] = fieldValue
	}
	return result
}

// GenerateJWT adds a unique jti claim so tokens can be revoked
func GenerateJWT(user models.User) (string, error) {
	userMap := StructToMap(user)
	expirationTime := time.Now().Add(48 * time.Hour).Unix()

	claims := jwt.MapClaims{
		"user": userMap,
		"exp":  expirationTime,
		"iat":  time.Now().Unix(),
		"jti":  uuid.NewString(), // unique token ID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	fmt.Println("Generated JWT for user:", user.Name, "| email:", user.Email)
	fmt.Println("Token is:", tokenString)
	return tokenString, nil
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure HMAC signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println("Error parsing JWT:", err)
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}
