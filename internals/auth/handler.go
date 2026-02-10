package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/models"
	"github.com/izzy-Ti/RaGO/internals/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type RegisterRequest struct {
	Name     string `json: "Name"`
	Email    string `json: "Email"`
	Password string `json: "Password"`
	Avatar   string `json: "avatar"`
}
type UpdateProfileRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Avatar *string `json:"avatar"`
}
type LoginRequest struct {
	Email    string `json: "Email"`
	Password string `json: "Password"`
}
type otpREQ struct {
	Otp      string `json:"otp"`
	Password string `json:"password"`
}
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := utils.ParseJSON(r, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var userModel models.User
	UserEmail := db.DB.Where("email = ?", req.Email).First(&userModel).Error
	if UserEmail == nil {
		utils.WriteJson(w, http.StatusUnauthorized, "Email already exists")
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := db.DB.Create(&user).Error; err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: utils.IsProd(),
		SameSite: func() http.SameSite {
			if utils.IsProd() {
				return http.SameSiteStrictMode
			}
			return http.SameSiteLaxMode
		}(),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if err := utils.SendWelcomeEmail(req.Email, req.Name, tokenString); err != nil {
		log.Println("Email send failed:", err)
		return
	}

	res := "login Successfully"
	resp := Response{
		Message: res,
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, resp)

}
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := utils.ParseJSON(r, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: utils.IsProd(),
		SameSite: func() http.SameSite {
			if utils.IsProd() {
				return http.SameSiteStrictMode
			}
			return http.SameSiteLaxMode
		}(),
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	res := "login Successfully"
	resp := Response{
		Message: res,
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, resp)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	res := Response{
		Message: "Logout successful",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)
}
func SendVerifyOTP(w http.ResponseWriter, r *http.Request) {
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()
	token, err := r.Cookie("token")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := utils.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user.VerifyOTP = utils.GenerateOTP()
	user.OTPExpireAt = int64(expiresAt)
	db.DB.Save(user)

	utils.SendOTPMail(user.Email, user.Name, user.VerifyOTP)
	res := Response{
		Message: "OTP sent successfully",
		Success: true,
	}
	utils.WriteJson(w, http.StatusOK, res)
}
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req otpREQ
	err := utils.ParseJSON(r, &req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	token, err := r.Cookie("token")
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := utils.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	if user.VerifyOTP != req.Otp || user.VerifyOTP == "" {
		utils.WriteJson(w, http.StatusUnauthorized, "invalid otp")
		return
	}
	if user.OTPExpireAt < time.Now().UnixMilli() {
		utils.WriteJson(w, http.StatusUnauthorized, "OTP expired")
		return
	}
	user.VerifyOTP = ""
	user.OTPExpireAt = 0
	user.IsAccVerified = true
	db.DB.Save(user)
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Welcome! Your account has been verified")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Thank you. Your account has been successfully verified.</p>
		<p>You can now sign in and start using your account.</p>
		<p>If you did not perform this action, you can reply to this email directly.</p>
	`, user.Name))

	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	d.DialAndSend(m)
	res := Response{
		Message: "OTP sent successfully",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)
}
func AuthUser(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*models.User)
	utils.WriteJson(w, http.StatusOK, user)
}
func SendResetOTP(w http.ResponseWriter, r *http.Request) {
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()
	token, _ := r.Cookie("token")
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := utils.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	email := user.Email
	user.ResetOTP = utils.GenerateOTP()
	user.ResetOTPExpireAt = int64(expiresAt)

	db.DB.Save(user)

	utils.SendOTPMail(email, user.Name, user.ResetOTP)
	res := Response{
		Message: "otp sent successfully",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)

}
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req otpREQ

	err := utils.ParseJSON(r, &req)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	token, _ := r.Cookie("token")
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user, err := utils.GetUserByID(userId)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	if user.ResetOTP != req.Otp || user.ResetOTP == "" {
		res := Response{
			Message: "invalid otp",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	if user.ResetOTPExpireAt < time.Now().UnixMilli() {
		res := Response{
			Message: "otp expired",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}
	user.Password = string(hashedPassword)
	db.DB.Save(user)
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("EMAIL"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Your password has been reset successfully")
	m.SetBody("text/html", fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your password has been successfully reset.</p>
		<p>If you did not perform this action, you can reply to this email directly.</p>
	`, user.Name))

	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAIL"), os.Getenv("PASSWORD"))
	d.DialAndSend(m)

	res := Response{
		Message: "password changed successfully",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)
}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req UpdateProfileRequest
	utils.ParseJSON(r, &req)
	user, _ := r.Context().Value("user").(*models.User)
	if req.Email != nil && *req.Email != "" {
		user.Email = *req.Email
	}
	if req.Name != nil && *req.Name != "" {
		user.Name = *req.Name
	}
	if req.Avatar != nil && *req.Avatar != "" {
		user.Avater = *req.Avatar
	}
	db.DB.Save(user)
	res := Response{
		Message: "profile changed successfully",
		Success: true,
	}

	utils.WriteJson(w, http.StatusOK, res)

}
func GoogleAuth(w http.ResponseWriter, r *http.Request) {}
