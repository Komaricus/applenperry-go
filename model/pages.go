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

type PagesAndFilesTable struct{}

func (PagesAndFilesTable) TableName() string {
	return "dbo.pages_and_files"
}

type PageAndFile struct {
	PagesAndFilesTable
	ID     string `gorm:"primarykey"`
	PageID string
	FileID string

	Page Page `gorm:"foreignKey:PageID"`
}
