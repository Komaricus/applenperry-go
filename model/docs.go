package model

import "time"

type DocsTable struct{}

func (DocsTable) TableName() string {
	return "dbo.docs"
}

type Document struct {
	DocsTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
