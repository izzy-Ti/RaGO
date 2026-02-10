package Middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/izzy-Ti/_server_setup/auth/internals/utils"
)

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

type Res struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			resp := Res{
				Message: "Unauthorized please login",
				Success: false,
			}
			utils.WriteJson(w, http.StatusUnauthorized, resp)
			return
		}
		userID, err := utils.UserId(token.Value, []byte(jwtSecret))
		if err != nil {
			resp := Res{
				Message: "Unauthorized please login",
				Success: false,
			}
			utils.WriteJson(w, http.StatusUnauthorized, resp)
			return
		}
		user, err := utils.GetUserByID(userID)
		if err != nil {
			resp := Res{
				Message: "Unauthorized please login",
				Success: false,
			}
			utils.WriteJson(w, http.StatusUnauthorized, resp)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
