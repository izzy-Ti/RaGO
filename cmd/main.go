package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/izzy-Ti/RaGO/internals/admin"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/embeddings"
	"github.com/izzy-Ti/RaGO/internals/server"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	_ = godotenv.Load()
	db.Connect()
	db.Migrate()
	DB, err := db.ConnectAstra()
	if err != nil {
		fmt.Errorf("failed to create vector collection: %s", err)
	}
	admin.ASTRA = DB
	embeddings.AS = DB
	db.CreateVectorCollection(DB)
	//embeddings.EmbedSite("https://israelashenafi.com/")

	handler := server.New()
	server.AuthRoutes(handler)
	server.AdminRoutes(handler)
	server.RagRoutes(handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://ra-go-frontend.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})
	handlerWithCORS := corsHandler.Handler(handler)
	log.Println("listening on" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), handlerWithCORS))
}
