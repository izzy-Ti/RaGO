package models

import "time"

type Chat struct {
	ID        uint `gorem:"primaryKey"`
	UserID    uint `gorem:"index;not null"`
	Title     string
	createdAt time.Time
	updatedAt time.Time

	Message []Message `gorm:""foreignKey:ChatID"`
}
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	ChatID    uint   `gorem:"index;not null"`
	Content   string `gorem:"type:text;not null"`
	Role      string `gorem:"type:varchar(20);not null"`
	createdAt time.Time
}
