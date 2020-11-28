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
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Word struct {
	AboutCiderTable
	ID   string `json:"id" gorm:"primarykey"`
	Name string `json:"text"`
}

type CiderAndFilesTable struct{}

func (CiderAndFilesTable) TableName() string {
	return "dbo.cider_and_files"
}

type CiderAndFile struct {
	CiderAndFilesTable
	ID      string `gorm:"primarykey"`
	CiderID string
	FileID  string

	Cider AboutCider `gorm:"foreignKey:CiderID"`
}
