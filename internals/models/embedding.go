package models

type KnowledgeChunk struct {
	ID      string    `json:"_id"`
	Content string    `json:"content"`
	Vector  []float32 `json:"$vector,omitempty"`
	Source  string    `json:"source"`
}
