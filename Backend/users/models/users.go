package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/mail"
	"os"
	"runtime"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"

	
	passwordvalidator "github.com/wagslane/go-password-validator"

	
	"github.com/alexedwards/argon2id"

	
	"gorm.io/gorm"
)

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
}

type UserModel struct {
	db *gorm.DB
}

type User struct {
	UserID   int    `gorm:"column:userid;primaryKey" json:"userid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	OTPCode  string `json:"-"`        
	Verified bool   `json:"verified"` 
}

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
		Verified: false,
	}

	res := e.db.Create(&user)
	if res.Error != nil {
		fmt.Println("error inserting user:", res.Error)
		return res.Error
	}

	otpCode := generateOTPCode()

	err = e.db.Model(&user).Updates(User{OTPCode: otpCode}).Error
	if err != nil {
		return err
	}


	err = sendOTPEmail(email, otpCode)
	if err != nil {
		fmt.Println("error sending OTP:", err)
		return err
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


func (e UserModel) UpdateVerificationStatus(user *User) error {
	return e.db.Model(user).Updates(map[string]interface{}{
		"otp_code": user.OTPCode,
		"verified": user.Verified,
	}).Error
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

func generateOTPCode() string {
	digits := "0123456789"
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		// fallback
		return "000000"
	}
	for i := 0; i < 6; i++ {
		buf[i] = digits[buf[i]%10]
	}
	return string(buf)
}

func sendOTPEmail(toEmail, otpCode string) error {
	from := sgmail.NewEmail("UniBazaar Support", "unibazaar.marketplace@gmail.com")
	subject := "Your UniBazaar OTP Code"
	to := sgmail.NewEmail("New User", toEmail)

	plainTextContent := fmt.Sprintf("Your OTP code is %s.\nIt expires in 5 minutes.", otpCode)
	htmlContent := fmt.Sprintf("<strong>Your OTP code is %s</strong><br>It expires in 5 minutes.", otpCode)

	message := sgmail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send OTP email:", err)
		return err
	}

	log.Printf("OTP email sent. Status Code: %d\n", response.StatusCode)
	return nil
}
