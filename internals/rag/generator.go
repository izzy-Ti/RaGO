package rag

import (
	"bytes"
	//"crypto/des"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	//"github.com/izzy-Ti/RaGO/internals/utils"
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
	SystemPrompt := `Decide whether this query requires database retrieval. Return ONLY valid JSON.`
	body := map[string]interface{}{
		"model": "llama-3.1-8b-instant",
		"messages": []map[string]string{
			{"role": "system",
				"content": SystemPrompt,
			},
			{
				"role":    "user",
				"content": query,
			},
		},
		"temperature": 0,
		"response_format": map[string]interface{}{
			"type": "json_schema",
			"json_schema": map[string]interface{}{
				"name": "search desicion",
				"schema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"answer": map[string]interface{}{
							"type":        "string",
							"description": "Direct answer if retrieval is not needed. Empty if search is required.",
						},
						"search": map[string]interface{}{
							"type":        "boolean",
							"description": "True if database search is required. False otherwise.",
						},
					},
					"required":             []string{"answer", "search"},
					"additionalProperties": false,
				},
			},
		},
	}
	json_body, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.groq.com/openai/v1/chat/completions",
		bytes.NewBuffer(json_body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("Groq API error %d: %s", res.StatusCode, string(body))
	}

	type SearchDecision struct {
		Ans    string `json:"answer"`
		Search bool   `json:"search"`
	}
	type GroqResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	var groqResp GroqResponse
	if err := json.NewDecoder(res.Body).Decode(&groqResp); err != nil {
		return "", fmt.Errorf("failed to decode Groq response: %v", err)
	}
	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("Groq API returned no choices")
	}
	content := groqResp.Choices[0].Message.Content
	fmt.Println("Groq returned:", content)

	var decision SearchDecision

	if err := json.Unmarshal([]byte(content), &decision); err != nil {
		return "", fmt.Errorf("invalid JSON from Groq: %v", err)
	}

	if decision.Search {
		contexts, err := Retriver(query)
		if err != nil {
			return "", err
		}
		return Generator(query, contexts)
	} else {
		return decision.Ans, nil
	}
}
