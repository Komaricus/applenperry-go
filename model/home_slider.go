package model

import "time"

type HomeSliderTable struct{}

func (HomeSliderTable) TableName() string {
	return "dbo.home_slider"
}

type HomeSliderItem struct {
	HomeSliderTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	Priority  int       `json:"priority"`
	FileID    string    `json:"fileId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	File File `json:"image" gorm:"foreignKey:FileID"`
}

type Slide struct {
	HomeSliderTable
	ID     string `json:"id" gorm:"primarykey"`
	FileID string `json:"file_id"`

	File File `json:"image" gorm:"foreignKey:FileID"`
}
