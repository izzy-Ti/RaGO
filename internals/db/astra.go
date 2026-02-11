package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	astradb "github.com/datastax/astra-db-go"
	"github.com/datastax/astra-db-go/options"
	"github.com/izzy-Ti/RaGO/internals/models"
)

var As *astradb.Db

func ConnectAstra() (*astradb.Db, error) {
	client := astradb.NewClient(options.WithToken(os.Getenv("ASTRA")))
	log.Println("Astradb connected")
	return client.Database(os.Getenv("ASTRA_END_POINT")), nil
}
func SaveData(As *astradb.Db, data models.KnowledgeChunk) {
	ctx := context.Background()
	collection := As.Collection("GoRag")

	collection.InsertOne(ctx, data)
}
func CreateVectorCollection(As *astradb.Db) error {
	ctx := context.Background()
	name := "GORag3"

	opts := &options.CollectionOptions{
		Vector: &options.VectorOptions{
			Dimension: 1024,
			Metric:    "cosine",
		},
	}
	_, err := As.CreateCollection(ctx, name, opts)
	if err != nil {
		if strings.Contains(err.Error(), "Collection already exists") {
			log.Printf("Collection %s already exists, skipping", name)
			return nil
		}
		return fmt.Errorf("failed to create collection: %w", err)
	}
	log.Printf("Created vector collection: %s (dim=%d metric=%s)\n", name, 1024, "COSINE")
	return nil
}
