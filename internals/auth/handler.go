package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/models"
	utils "github.com/izzy-Ti/RaGO/internals/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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
type GoogleRequest struct {
	Token string `json:"token"`
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
	if req.Email == "" || req.Name == "" || req.Password == "" {
		res := Response{
			Message: "Missing data",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	UserEmail := db.DB.Where("email = ?", req.Email).First(&userModel).Error
	if UserEmail == nil {
		res := Response{
			Message: "Email already exists",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		res := Response{
			Message: "bcrypting error",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
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
		HttpOnly: true,
		Secure:   true,
		SameSite:  http.SameSiteNoneMode,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	verificationUrl := ""
	subject := "Password reset successfull"
	html := fmt.Sprintf(`
        <p>Hi %s,</p>
        <p>Welcome! Please verify your email by clicking the link below:</p>
        <a href="%s">Verify Email</a>
        <p>If this wasn’t you, please ignore this email.</p>
    `, user.Name, verificationUrl)
	er := utils.Sendemail(user.Email, user.Name, subject, html)
	if er != nil {
		utils.WriteError(w, http.StatusUnauthorized, er)
	}

	res := "Registration Successfully"
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
		res := Response{
			Message: "can't find the email",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		res := Response{
			Message: "Incorrcet credentials",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
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
		HttpOnly: true,
		Secure:   true,
		SameSite:  http.SameSiteNoneMode,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
	})
	subject := "Login alert"
	html := `<h2 style="color:#2e7d32;">Login Successful</h2>
			<p>You have successfully logged into your account.</p>
			<p>If this was you, no further action is required.</p>
			<p>If you do not recognize this login activity, please reset your password immediately to secure your account.</p>
			<hr style="margin-top:20px;">
			<p style="font-size:12px;color:#888;">
			This is an automated message. Please do not reply to this email.
			</p>
			`
	er := utils.Sendemail(req.Email, user.Name, subject, html)
	if er != nil {
		utils.WriteError(w, http.StatusUnauthorized, er)
	}
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
		res := Response{
			Message: "Please login",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	userId, err := utils.UserId(token.Value, []byte(jwtSecret))
	if err != nil {
		res := Response{
			Message: "user not found",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	user, err := utils.GetUserByID(userId)
	if err != nil {
		res := Response{
			Message: "user not found",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	user.VerifyOTP = utils.GenerateOTP()
	user.OTPExpireAt = int64(expiresAt)
	db.DB.Save(user)

	subject := "OTP verfication"
	html := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your one-time verification code is:</p>
		<h2 style="letter-spacing:2px;">%s</h2>
		<p>This code will expire soon. Do not share it with anyone.</p>
		<p>If you didn’t request this, you can ignore this email.</p>
	`, user.Name, user.VerifyOTP)

	er := utils.Sendemail(user.Email, user.Name, subject, html)
	if er != nil {
		utils.WriteError(w, http.StatusUnauthorized, er)
	}
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
	user.ResetOTP = utils.GenerateOTP()
	user.ResetOTPExpireAt = int64(expiresAt)

	db.DB.Save(user)

	subject := "OTP verfication"
	html := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your one-time verification code is:</p>
		<h2 style="letter-spacing:2px;">%s</h2>
		<p>This code will expire soon. Do not share it with anyone.</p>
		<p>If you didn’t request this, you can ignore this email.</p>
	`, user.Name, user.ResetOTP)

	er := utils.Sendemail(user.Email, user.Name, subject, html)
	if er != nil {
		utils.WriteError(w, http.StatusUnauthorized, er)
	}
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
	subject := "Password reset successfull"
	html := `
		<p>Hi %s,</p>
		<p>Your password has been successfully reset.</p>
		<p>If you made this change, you can now log in with your new password.</p>
		<p>If you did not request this password reset, please contact support immediately to secure your account.</p>
	`

	er := utils.Sendemail(user.Email, user.Name, subject, html)
	if er != nil {
		utils.WriteError(w, http.StatusUnauthorized, er)
	}

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
func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req GoogleRequest
	utils.ParseJSON(r, &req)
	if req.Token == "" {
		res := Response{
			Message: "Missing data",
			Success: false,
		}
		utils.WriteJson(w, http.StatusUnauthorized, res)
		return
	}
	payload, err := idtoken.Validate(context.Background(), req.Token, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		http.Error(w, "invalid google token", http.StatusUnauthorized)
		return
	}
	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	sub := payload.Claims["sub"].(string)

	var user models.User

	err = db.DB.Where("email=?", email).First(&user).Error
	if err != nil {
		user = models.User{
			Name:          name,
			Email:         email,
			Avater:        picture,
			GoogleId:      sub,
			AuthType:      "google",
			IsAccVerified: true,
		}
		db.DB.Create(&user)
	} else {
		if !user.IsAccVerified {
			user.Name = name
			user.Avater = picture
			user.GoogleId = sub
			user.AuthType = "google"
			user.IsAccVerified = true
			db.DB.Save(&user)
		}
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
