package models

import "time"

type Chat struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index;not null"`
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time

	Messages []Message `gorm:"foreignKey:ChatID"`
}
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	ChatID    uint   `gorm:"index;not null"`
	Content   string `gorm:"type:text;not null"`
	Role      string `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
}
