package rag

import (
	"context"
	"fmt"

	"github.com/datastax/astra-db-go/filter"
	"github.com/datastax/astra-db-go/options"
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
	cursor := col.Find(ctx, filter.F{},
		options.WithCollectionLimit(3),
		options.WithCollectionSort(map[string]interface{}{"$vector": queryVec}),
		options.WithCollectionIncludeSimilarity(true),
	)
	if err := cursor.Err(); err != nil {
		fmt.Println("CRITICAL: Search Error:", err)
	}
	defer cursor.Close(ctx)
	total, _ := col.CountDocuments(ctx, map[string]interface{}{}, 0)
	fmt.Println("Total documents in collection:", total)
	var results []string
	count := 0
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			fmt.Println("Decode error:", err)
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
