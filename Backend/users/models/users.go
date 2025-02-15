package models

import (
	"fmt"
	"net/mail"
	"runtime"
	"strings"

	passwordvalidator "github.com/wagslane/go-password-validator"

	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
)

// for now we are only considering local universities.
var validEDUDomains = map[string]bool{
	"ufl.edu":         true,
	"fsu.edu":         true,
	"ucf.edu":         true,
	"usf.edu":         true,
	"fiu.edu":         true,
	"fau.edu":         true,
	"fgcu.edu":        true,
	"unf.edu":         true,
	"famu.edu":        true,
	"ncf.edu":         true,
	"floridapoly.edu": true,
	"uwf.edu":         true,
}

type UserModel struct {
	db *gorm.DB
}

type User struct {
	UserID   int    `gorm:"column:userid;primaryKey" json:"userid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// increase iterations to improve security, can also increase key length, saltlength should not be changed 16 is recommended
var params = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

const minEntropyBits = 60

func ValidatePassword(password string) error {
	err := passwordvalidator.Validate(password, minEntropyBits)
	if err != nil {
		return fmt.Errorf("password is too weak: %v", err)
	}
	return nil
}

func (e UserModel) Insert(id int, name string, email string, password string) error {
	if err := ValidateEduEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("error hashing password:", err)
		return err
	}

	user := User{
		UserID:   id,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	res := e.db.Create(&user)
	if res.Error != nil {
		fmt.Println("error inserting user:", res.Error)
		return res.Error
	}

	return nil
}

func (e UserModel) GetUserIdByEmail(email string) (int, error) {
	var user User
	res := e.db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		fmt.Println("error while reading user information")
		return 0, res.Error
	}

	return user.UserID, nil
}

func (e UserModel) Update(email string, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return err
	}

	res := e.db.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword)
	if res.Error != nil {
		fmt.Println("Error while updating password")
		return res.Error
	}
	return nil
}

func (e UserModel) Delete(email string) error {
	res := e.db.Where("email = ?", email).Delete(&User{})
	if res.Error != nil {
		fmt.Println("error while deleting user")
		return res.Error
	}
	return nil

}

func (e UserModel) Read(email string) (*User, error) {
	var user User
	res := e.db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		fmt.Println("error while reading user information")
		return nil, res.Error
	}
	return &user, nil

}

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func ValidateEduEmail(email string) error {
	addr, err := mail.ParseAddress(strings.TrimSpace(email))
	if err != nil {
		return fmt.Errorf("invalid email address: %v", err)
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email address: missing @ or domain")
	}

	domain := strings.ToLower(parts[1])

	if !strings.HasSuffix(domain, ".edu") {
		return fmt.Errorf("invalid email: must be a .edu address")
	}

	if !validEDUDomains[domain] {
		return fmt.Errorf("unrecognized .edu domain: %s", domain)
	}

	return nil
}
