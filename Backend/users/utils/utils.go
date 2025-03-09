package utils

import (
	"fmt"
	"os"
	"reflect"
	"time"
	"users/models"

	"github.com/golang-jwt/jwt/v5"
)

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

func GenerateJWT(user models.User) (string, error) {
	userMap := StructToMap(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":   userMap,
		"expiry": time.Now().AddDate(0, 0, 2),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	fmt.Println("Generated JWT for user " + user.Name + " and email: " + user.Email + "is: " + tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return token, nil
}
