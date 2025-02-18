package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	FullName  string    `gorm:"type:varchar(255)" json:"full_name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Type      int       `gorm:"not null" json:"type"`                           // type merujuk ke ID pada TypeUser
	TypeInfo  TypeUser  `gorm:"foreignKey:Type;references:ID" json:"type_info"` // foreign key ke kolom Type pada User yang merujuk ke ID pada TypeUser
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TypeUser struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Type string `gorm:"type:varchar(30);not null" json:"type"`
}

func (User) TableName() string {
	return "user"
}

func (TypeUser) TableName() string {
	return "type_user"
}