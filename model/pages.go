package model

import "time"

type PagesTable struct{}

func (PagesTable) TableName() string {
	return "dbo.pages"
}

type Page struct {
	PagesTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
