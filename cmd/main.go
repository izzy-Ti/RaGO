package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db.Connect()
	db.Migrate()
	DB, err := db.ConnectAstra()
	if err != nil {
		fmt.Errorf("failed to create vector collection: %s", err)
	}
	db.CreateVectorCollection(DB)

	handler := server.New()
	server.AuthRoutes(handler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
