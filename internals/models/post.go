package models

type Posts struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"not null"`
	Uploadby  User   `gorm:"foreignKey:UploadedBy;references:ID"`
	CreatedAt uint64 `gorm:"not null;autoCreateTime"`
}
