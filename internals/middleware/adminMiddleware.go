package middleware

import (
	"context"
	"net/http"

	utils "github.com/izzy-Ti/RaGO/internals/Utils"
)

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, "no Token login pls")
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
		if user.Role != "ADMIN" {
			res := Response{
				Message: "Sorry you are not an admin",
				Success: false,
			}
			utils.WriteJson(w, http.StatusUnauthorized, res)
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
