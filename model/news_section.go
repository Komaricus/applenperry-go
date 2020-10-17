package model

import "time"

type NewsSectionsTable struct{}

func (NewsSectionsTable) TableName() string {
	return "dbo.news_sections"
}

type NewsSection struct {
	NewsSectionsTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
