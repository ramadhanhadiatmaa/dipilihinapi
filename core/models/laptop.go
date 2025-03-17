package models

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

// Laptop adalah representasi tabel laptop
type Laptop struct {
	ID        int    `json:"id" gorm:"type:int4;primaryKey"`
	Name      string `json:"name" gorm:"type:varchar(512)"`
	Brand     string `json:"brand" gorm:"type:varchar(512)"`
	Cpu       string `json:"cpu" gorm:"type:varchar(512)"`
	Ram       string `json:"ram" gorm:"type:varchar(512)"`
	Storage   string `json:"storage" gorm:"type:varchar(512)"`
	Gpu       string `json:"gpu" gorm:"type:varchar(512)"`
	Display   string `json:"display" gorm:"type:varchar(512)"`
	Weight    string `json:"weight" gorm:"type:varchar(512)"`
	Price     int    `json:"price" gorm:"type:int4"`
	ImageUrl  string `json:"image_url" gorm:"type:varchar(512)"`
	TokpedUrl string `json:"tokped_url" gorm:"type:varchar(512)"`
	TiktokUrl string `json:"tiktok_url" gorm:"type:varchar(512)"`
	ShopeeUrl string `json:"shopee_url" gorm:"type:varchar(512)"`
}

// TableName mengubah nama tabel default menjadi "laptop"
func (Laptop) TableName() string {
	return "laptop"
}

// GetLaptopsByIds mengambil data laptop berdasarkan list ID
func GetLaptopsByIds(ids []int) ([]Laptop, error) {
	var laptops []Laptop
	result := DB.Where("id IN ?", ids).Find(&laptops)
	return laptops, result.Error
}
