package rag

import (
	"context"

	"github.com/izzy-Ti/RaGO/internals/embeddings"
)

func Retriver(query string) ([]string, error) {
	ctx := context.Background()
	queryVec, err := embeddings.EmbedText(query)
	if err != nil {
		return nil, err
	}
	col := embeddings.AS.Collection("GORag3")
	body := map[string]interface{}{
		"find": map[string]interface{}{
			"sort": map[string]interface{}{
				"$vector": queryVec,
			},
			"options": map[string]interface{}{
				"limit":             3,
				"includeSimilarity": true,
			},
		},
	}

	cursor := col.Find(ctx, body)
	defer cursor.Close(ctx)

	var results []string
	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		if content, ok := doc["content"].(string); ok {
			results = append(results, content)
		}
	}
	return results, nil
}
