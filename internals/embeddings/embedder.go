package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	astradb "github.com/datastax/astra-db-go"
)

var AS *astradb.Db

func EmbedText(text string) ([]float32, error) {
	body := map[string]interface{}{
		"model": "voyage-3-large",
		"input": []string{text},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(
		"POST",
		"https://api.voyageai.com/v1/embeddings",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("VOYAGE_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("voyage error: %s", string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	var vectors [][]float32
	for _, d := range result.Data {
		vectors = append(vectors, d.Embedding)
	}

	return result.Data[0].Embedding, nil
}
func FetchRenderedHTML(url string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.InnerHTML("body", &html, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}

	return html, nil
}
func RemoveSpecialChars(text string) string {
	re := regexp.MustCompile(`[^\P{So}\p{L}\p{N}\p{P}\p{Z}]`)
	return re.ReplaceAllString(text, "")
}
func EmbedSite(url string) ([][]float32, error) {
	ctx := context.Background()
	html, _ := FetchRenderedHTML(url)
	//stream byte for editing the html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	doc.Find("script, style, noscript").Remove()
	text := strings.TrimSpace(doc.Text())
	text = strings.Join(strings.Fields(text), " ")
	text = RemoveSpecialChars(text)
	if text == "" {
		return nil, fmt.Errorf("empty site content")
	}

	chunks := ChunkText(text, 800)
	var vectors [][]float32

	for i, chunk := range chunks {
		vec, err := EmbedText(chunk)
		if err != nil {
			log.Printf("embedding failed for chunk %d: %v\nChunk preview: %.50s", i+1, err, chunk)
			continue
		}
		docs := map[string]interface{}{
			"content": chunk,
			"$vector": vec,
		}
		col := AS.Collection("GORag3")
		_, err = col.InsertOne(ctx, docs)
		if err != nil {
			log.Println("insert failed:", err)
			continue
		}
		log.Printf("chunk %d embedded (%d chars)", i+1, len(chunk))
		vectors = append(vectors, vec)
		time.Sleep(25 * time.Second)
	}
	return vectors, nil

}
func ChunkText(text string, size int) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += size {
		end := i + size

		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
