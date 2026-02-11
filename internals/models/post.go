package models

type Posts struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null"`
	Content   string `gorm:"not null"`
	Uploadby  uint
	User      User   `gorm:"foreignKey:Uploadby"`
	CreatedAt uint64 `gorm:"not null;autoCreateTime"`
}
