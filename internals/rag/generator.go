package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func Generator(query string, contexts []string) (string, error) {
	prompt := "You are part of the RAGO team. Speak in first person plural (use 'we', 'our', 'us'). " +
	"Do NOT say 'they' or refer to RAGO in third person. " +
	"Keep the answer concise (maximum 4–6 sentences). " +
	"Be confident, supportive, and visionary while staying clear and professional.\n\n" +
	"Use the following context to answer:\n\n" + strings.Join(contexts, "\n") + "\n\nQuestion: " + query + "\nAnswer:"

	body := map[string]interface{}{
		"model": "llama-3.1-8b-instant",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Groq API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no output from generator")
	}

	return result.Choices[0].Message.Content, nil
}
func RAG(query string) (string, error) {
	contexts, err := Retriver(query)
	if err != nil {
		return "", err
	}

	return Generator(query, contexts)
}
