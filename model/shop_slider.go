package model

import "time"

type ShopSliderTable struct{}

func (ShopSliderTable) TableName() string {
	return "dbo.shop_slider"
}

type ShopSliderItem struct {
	ShopSliderTable
	ID          string    `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Priority    int       `json:"priority"`
	FileID      string    `json:"fileId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	File File `json:"image" gorm:"foreignKey:FileID"`
}

type ShopSlide struct {
	ShopSliderTable
	ID          string `json:"id" gorm:"primarykey"`
	Name        string `json:"header"`
	Description string `json:"description"`
	Link        string `json:"link"`
	FileID      string `json:"fileId"`

	File File `json:"image" gorm:"foreignKey:FileID"`
}
