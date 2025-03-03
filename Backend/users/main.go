package main

import (
	"fmt"
	"users/server"

	"github.com/joho/godotenv"
)

type User2 struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func CreateUser(firstName string, lastName string, email string, password string) *User2 {
	user := User2{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	return &user
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	server.InitServer()
}
