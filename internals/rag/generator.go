package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Generator(query string, contexts []string) (string, error) {
	prompt := "Answer the question using the following context:\n\n" +
		strings.Join(contexts, "\n") +
		"\n\nQuestion: " + query + "\nAnswer:"

	body := map[string]interface{}{
		"model": "voyage-3-large",
		"input": prompt,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.voyageai.com/v1/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VOYAGE_API_KEY"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Output []struct {
			Text string `json:"text"`
		} `json:"output"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Output) == 0 {
		return "", fmt.Errorf("no output from generator")
	}

	return result.Output[0].Text, nil
}
func RAG(query string) (string, error) {
	contexts, err := Retriver(query)
	if err != nil {
		return "", err
	}

	return Generator(query, contexts)
}
