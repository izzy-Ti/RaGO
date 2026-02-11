package rag

import (
	"context"
	"fmt"

	"github.com/izzy-Ti/RaGO/internals/embeddings"
)

func Retriver(query string) ([]string, error) {
	ctx := context.Background()
	queryVec, err := embeddings.EmbedText(query)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Generated embedding of length: %d\n", len(queryVec))
	col := embeddings.AS.Collection("GORag3")
	filter := map[string]interface{}{
		"sort": map[string]interface{}{
			"$vector": queryVec,
		},
		"options": map[string]interface{}{
			"limit":             3,
			"includeSimilarity": true,
		},
	}

	cursor := col.Find(ctx, filter)
	defer cursor.Close(ctx)

	var results []string
	count := 0
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		if content, ok := doc["content"].(string); ok {
			results = append(results, content)
			count++
		}
	}
	fmt.Println("Retriever found documents:", count)
	return results, nil
}
