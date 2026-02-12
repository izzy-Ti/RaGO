package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izzy-Ti/RaGO/internals/admin"
	"github.com/izzy-Ti/RaGO/internals/auth"
	"github.com/izzy-Ti/RaGO/internals/chat"
	Middleware "github.com/izzy-Ti/RaGO/internals/middleware"
)

func AuthRoutes(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", auth.Register).Methods("POST")
	userRouter.HandleFunc("/login", auth.Login).Methods("POST")
	userRouter.HandleFunc("/logout", auth.Logout).Methods("POST")
	userRouter.HandleFunc("/sendotp", auth.SendVerifyOTP).Methods("POST")
	userRouter.HandleFunc("/verifyotp", auth.VerifyOTP).Methods("POST")
	userRouter.HandleFunc("/sendresetotp", auth.SendResetOTP).Methods("POST")
	userRouter.HandleFunc("/verifyreset", auth.ResetPassword).Methods("POST")
	userRouter.Handle("/auth", Middleware.IsAuth(http.HandlerFunc(auth.AuthUser))).Methods("POST")
	userRouter.Handle("/updateprofile", Middleware.IsAuth(http.HandlerFunc(auth.UpdateProfile))).Methods("POST")
}
func AdminRoutes(r *mux.Router) {
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Handle("/post", Middleware.AdminAuth(http.HandlerFunc(admin.SaveAdminPost))).Methods("POST")
}
func RagRoutes(r *mux.Router) {
	ragRouter := r.PathPrefix("/ask").Subrouter()
	ragRouter.Handle("/que", Middleware.IsAuth(http.HandlerFunc(chat.Ask)))
	//ragRouter.HandleFunc("/que", chat.Ask).Methods("POST")
}
