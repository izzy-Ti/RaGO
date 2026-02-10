package main

import (
	"log"
	"net/http"

	"github.com/izzy-Ti/RaGO/internals/auth"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db.Connect()
	db.Migrate()

	handler := server.New()
	auth.AuthRoutes(handler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
