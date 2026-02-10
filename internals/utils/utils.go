package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/izzy-Ti/_server_setup/auth/internals/db"
	"github.com/izzy-Ti/_server_setup/auth/internals/models"
	"gopkg.in/gomail.v2"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing req")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}
func IsProd() bool {
	return os.Getenv("APP_ENV") == "production"
}
func SendWelcomeEmail(to, name, token string) error {
	verificationUrl := fmt.Sprintf("https://fmbls.vercel.app/verifyemail")

	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome! Please verify your email")
	m.SetBody("text/html", fmt.Sprintf(`
        <p>Hi %s,</p>
        <p>Welcome! Please verify your email by clicking the link below:</p>
        <a href="%s">Verify Email</a>
        <p>If this wasn’t you, please ignore this email.</p>
    `, name, verificationUrl))

	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	return d.DialAndSend(m)
}
func UserId(tokenStr string, jwtSecret []byte) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		log.Println("Token parsing error:", err) // log the error for debugging
		return 0, errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("sub not found")
	}

	return uint(userID), nil
}
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(100000 + rand.Intn(900000))
}
func SendOTPMail(to, name, otp string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome! Please verify your email")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your one-time verification code is:</p>
		<h2>%s</h2>
		<p>This code will expire soon. Do not share it with anyone.</p>
		<p>If you didn’t request this, you can ignore this email.</p>
		`, name, otp))
	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	return d.DialAndSend(m)
}