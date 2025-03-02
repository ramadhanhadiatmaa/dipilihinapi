package models

type Status struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Status string `gorm:"type:varchar(30);not null" json:"status"`
}

func (Status) TableName() string {
	return "status"
}