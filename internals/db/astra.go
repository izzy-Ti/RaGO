package db

import (
	"context"
	"log"
	"os"

	astradb "github.com/datastax/astra-db-go"
	"github.com/datastax/astra-db-go/options"
	"github.com/izzy-Ti/RaGO/internals/models"
)

func ConnectAstra() (*astradb.Db, error) {
	client := astradb.NewClient(options.WithToken(os.Getenv("ASTRA")))
	log.Println("Astradb connected")
	return client.Database(os.Getenv("ASTRA_END_POINT")), nil
}
func SaveData(db *astradb.Db, data models.KnowledgeChunk) {
	ctx := context.Background()
	collection := db.Collection("GoRag")

	collection.InsertOne(ctx, data)
}
func CreateVectorCollection(db *astradb.Db) error {
	ctx := context.Background()

	opts := &options.CollectionOptions{
		Vector: &options.VectorOptions{
			Dimension: 1024,
			Metric:    "cosine",
		},
	}
	_, err := db.CreateCollection(ctx, "GORag", opts)
	if err != nil {
		log.Printf("error creating astra collection err= %s", err)
		return nil
	}
	log.Printf("Created vector collection: %s (dim=%d metric=%s)\n", "GORag", 1536, "COSINE")
	return nil
}
