package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
	passwordvalidator "github.com/wagslane/go-password-validator"
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

var params = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

const minEntropyBits = 60

// Only US phone is allowed: 10 digits, optional +1 prefix
var phoneRegex = regexp.MustCompile(`^\+?1?\d{10}$`)

// ValidatePhone checks US phone format
func ValidatePhone(phone string) error {
	trimmed := strings.ReplaceAll(phone, " ", "")
	trimmed = strings.ReplaceAll(trimmed, "-", "")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return nil
	}
	if !phoneRegex.MatchString(trimmed) {
		return errors.New("invalid US phone number; must be 10 digits with optional +1 prefix")
	}
	return nil
}

func ValidatePassword(password string) error {
	err := passwordvalidator.Validate(password, minEntropyBits)
	if err != nil {
		return fmt.Errorf("password is too weak: %v", err)
	}
	return nil
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
		return "000000"
	}
	for i := 0; i < 6; i++ {
		buf[i] = digits[buf[i]%10]
	}
	return string(buf)
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, params)
}

// User is our GORM model
type User struct {
	UserID              int    `gorm:"column:userid;primaryKey" json:"userid"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"-"`
	OTPCode             string `json:"-"`
	FailedResetAttempts int    `json:"-"`
	Verified            bool   `json:"-"`
	Phone               string `json:"phone"`
}

// UserModel now exports the DB field
type UserModel struct {
	DB *gorm.DB
}

// Insert creates a new user
func (e UserModel) Insert(id int, name, email, password, phone string) error {
	// Validate
	if err := ValidateEduEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	if err := ValidatePhone(phone); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	//TODO: Return the user to the handler otherwise we will not be able to keep track of the user id
	user := User{
		UserID:   id,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Verified: false,
		Phone:    phone,
	}

	// Create user
	if err := e.DB.Create(&user).Error; err != nil {
		return err
	}

	// Generate OTP
	otpCode := generateOTPCode()
	user.OTPCode = otpCode
	if err := e.DB.Save(&user).Error; err != nil {
		return err
	}

	// Send OTP email
	if err := sendOTPEmail(email, otpCode, "Your UniBazaar OTP Code"); err != nil {
		return err
	}
	return nil
}

// Update just the password
func (e UserModel) Update(email, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	res := e.DB.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword)
	return res.Error
}

// UpdateName changes the user's name
func (e UserModel) UpdateName(email, newName string) error {
	var user User
	if err := e.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.Name = newName
	return e.DB.Save(&user).Error
}

// UpdatePhone changes the user's phone
func (e UserModel) UpdatePhone(email, newPhone string) error {
	var user User
	if err := e.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	if err := ValidatePhone(newPhone); err != nil {
		return err
	}
	user.Phone = newPhone
	return e.DB.Save(&user).Error
}

// Delete removes a user by email
func (e UserModel) Delete(email string) error {
	res := e.DB.Where("email = ?", email).Delete(&User{})
	return res.Error
}

// Read fetches a user by email
func (e UserModel) Read(email string) (*User, error) {
	var user User
	res := e.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

// UpdateVerificationStatus marks user verified
func (e UserModel) UpdateVerificationStatus(user *User) error {
	return e.DB.Model(user).Updates(map[string]interface{}{
		"verified": user.Verified,
	}).Error
}

// SaveUser saves an existing user
func (e UserModel) SaveUser(user *User) error {
	return e.DB.Save(user).Error
}

// SendSecurityAlert example
func (e UserModel) SendSecurityAlert(user *User) error {
	return sendOTPEmail(user.Email, "Suspicious attempts detected", "UniBazaar Security Alert")
}

// GetUserIdByEmail fetches userID
func (e UserModel) GetUserIdByEmail(email string) (int, error) {
	var user User
	res := e.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return 0, res.Error
	}
	return user.UserID, nil
}

// InitiatePasswordReset sets an OTP for password reset
func (e UserModel) InitiatePasswordReset(email string) error {
	user, err := e.Read(email)
	if err != nil {
		return fmt.Errorf("user not found or DB error: %w", err)
	}
	user.FailedResetAttempts = 0
	if err := e.DB.Save(user).Error; err != nil {
		return err
	}
	otpCode := generateOTPCode()
	user.OTPCode = otpCode
	if err := e.DB.Save(user).Error; err != nil {
		return err
	}
	if err := sendOTPEmail(email, otpCode, "UniBazaar Password Reset Code"); err != nil {
		return err
	}
	return nil
}

func (e UserModel) VerifyResetCodeAndSetNewPassword(email, code, newPassword string) error {
	user, err := e.Read(email)
	if err != nil {
		return fmt.Errorf("user not found or DB error: %w", err)
	}
	if user.OTPCode != code {
		user.FailedResetAttempts++
		if user.FailedResetAttempts >= 3 {
			_ = sendOTPEmail(user.Email, "Suspicious attempts detected", "UniBazaar Security Alert")
			user.FailedResetAttempts = 0
			user.OTPCode = ""
		}
		if err := e.DB.Save(user).Error; err != nil {
			return err
		}
		return fmt.Errorf("invalid OTP code")
	}
	user.FailedResetAttempts = 0
	if err := ValidatePassword(newPassword); err != nil {
		return err
	}
	hashed, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	user.OTPCode = ""
	if err := e.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func sendOTPEmail(toEmail, code, subject string) error {
	from := sgmail.NewEmail("UniBazaar Support", "unibazaar.marketplace@gmail.com")
	to := sgmail.NewEmail("User", toEmail)
	plainText := fmt.Sprintf("Your code is %s.\nIt expires in 5 minutes.", code)
	htmlContent := fmt.Sprintf("<strong>Your code is %s</strong><br>It expires in 5 minutes.", code)
	message := sgmail.NewSingleEmail(from, subject, to, plainText, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}
	log.Printf("Email sent. Status Code: %d\n", response.StatusCode)
	return nil
}

func CreateUser(name string, email string, phone string) *User {
	user := User{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	return &user
}
