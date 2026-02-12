package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/izzy-Ti/RaGO/internals/admin"
	"github.com/izzy-Ti/RaGO/internals/db"
	"github.com/izzy-Ti/RaGO/internals/embeddings"
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
	admin.ASTRA = DB
	embeddings.AS = DB
	db.CreateVectorCollection(DB)
	//embeddings.EmbedSite("https://israelashenafi.com/")

	handler := server.New()
	server.AuthRoutes(handler)
	server.AdminRoutes(handler)
	server.RagRoutes(handler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":10000 ", handler))
}
