package models

type RagKnowladge struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
}
