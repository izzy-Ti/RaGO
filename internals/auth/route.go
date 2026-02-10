package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izzy-Ti/_server_setup/auth/internals/Middleware"
)

func AuthRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", Register).Methods("POST")
	userRouter.HandleFunc("/login", Login).Methods("POST")
	userRouter.HandleFunc("/logout", Logout).Methods("POST")
	userRouter.HandleFunc("/sendotp", SendVerifyOTP).Methods("POST")
	userRouter.HandleFunc("/verifyotp", VerifyOTP).Methods("POST")
	userRouter.HandleFunc("/sendresetotp", SendResetOTP).Methods("POST")
	userRouter.HandleFunc("/verifyreset", ResetPassword).Methods("POST")
	userRouter.Handle("/auth", Middleware.IsAuth(http.HandlerFunc(AuthUser))).Methods("POST")
	userRouter.Handle("/updateprofile", Middleware.IsAuth(http.HandlerFunc(UpdateProfile))).Methods("POST")
}
