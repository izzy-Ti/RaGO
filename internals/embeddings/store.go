package embeddings

type Posts struct {
	ID      string    `json:"_id,omitempty"`
	PostID  uint      `json:"post_id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Vector  []float32 `json:"$vector"`
}
