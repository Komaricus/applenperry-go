package model

import "time"

type AboutCiderTable struct{}

func (AboutCiderTable) TableName() string {
	return "dbo.about_cider"
}

type AboutCider struct {
	AboutCiderTable
	ID          string    `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Size        int       `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Word struct {
	AboutCiderTable
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"text"`
	Size int    `json:"weight"`
}
