package models

import "time"

type Chat struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(100);not null" json:"email"`
	Message   string    `gorm:"type:varchar(1000);not null" json:"message"`
	Status    int       `gorm:"not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Chat) TableName() string {
	return "chat"
}
