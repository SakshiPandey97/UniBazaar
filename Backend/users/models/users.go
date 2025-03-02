package models

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/mail"
	"os"
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

type UserModel struct {
	db *gorm.DB
}

type User struct {
	UserID              int    `gorm:"column:userid;primaryKey" json:"userid"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Password            string `json:"-"`
	OTPCode             string `json:"-"`
	FailedResetAttempts int    `json:"-"`
	Verified            bool   `json:"-"`
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

func (e UserModel) Insert(id int, name, email, password string) error {
	if err := ValidateEduEmail(email); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user := User{
		UserID:   id,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Verified: false,
	}
	fmt.Println("going to add user to the db")
	if err := e.db.Create(&user).Error; err != nil {
		fmt.Println("error occurred while inserting into the db")
		return err
	}
	fmt.Println("added user to db")
	otpCode := generateOTPCode()
	if err := e.db.Model(&user).Updates(User{OTPCode: otpCode}).Error; err != nil {
		return err
	}
	if err := sendOTPEmail(email, otpCode, "Your UniBazaar OTP Code"); err != nil {
		return err
	}
	return nil
}

func (e UserModel) Update(email, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	res := e.db.Model(&User{}).Where("email = ?", email).Update("password", hashedPassword)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (e UserModel) Delete(email string) error {
	res := e.db.Where("email = ?", email).Delete(&User{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (e UserModel) Read(email string) (*User, error) {
	var user User
	res := e.db.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (e UserModel) UpdateVerificationStatus(user *User) error {
	return e.db.Model(user).Updates(map[string]interface{}{
		"verified": user.Verified,
	}).Error
}

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, params)
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

// generateOTPCode creates a random 6-digit numeric code.
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

func (e UserModel) GetUserIdByEmail(email string) (int, error) {
	var user User
	res := e.db.Where("email = ?", email).First(&user)

	if res.Error != nil {
		fmt.Println("error while reading user information")
		return 0, res.Error
	}

	return user.UserID, nil
}

// sendOTPEmail sends the 6-digit code via SendGrid.
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

func (e UserModel) InitiatePasswordReset(email string) error {
	user, err := e.Read(email)
	if err != nil {
		return fmt.Errorf("user not found or DB error: %w", err)
	}
	user.FailedResetAttempts = 0
	if err := e.db.Save(user).Error; err != nil {
		return err
	}
	otpCode := generateOTPCode()
	user.OTPCode = otpCode
	if err := e.db.Save(user).Error; err != nil {
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
		fmt.Println(user)
		if user.FailedResetAttempts >= 3 {
			_ = sendOTPEmail(user.Email, "Suspicious attempts detected", "UniBazaar Security Alert")
			user.FailedResetAttempts = 0
			user.OTPCode = ""
		}
		if err := e.db.Save(user).Error; err != nil {
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
	if err := e.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Optional helper if needed in routes:
func (e UserModel) SaveUser(user *User) error {
	return e.db.Save(user).Error
}

func (e UserModel) SendSecurityAlert(user *User) error {
	return sendOTPEmail(user.Email, "Suspicious attempts detected", "UniBazaar Security Alert")
}
