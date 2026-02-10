package models

type Embedding struct {
	ID         uint64    `gorm:"primaryKey"`
	DocumentID uint64    `gorm:"not null"`
	Content    string    `gorm:"type:text"`
	Vector     []float32 `gorm:"type:vector(1536)"`
	CreatedAt  uint64    `gorm:"autoCreateTime"`
}
