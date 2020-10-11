package model

import "time"

type NewsTable struct{}

func (NewsTable) TableName() string {
	return "dbo.news"
}

type News struct {
	NewsTable
	ID          string      `json:"id" gorm:"primarykey"`
	Name        string      `json:"name"`
	SectionID   string      `json:"sectionId"`
	Subheader   string      `json:"subheader"`
	Description string      `json:"description"`
	FileID      string      `json:"fileId"`
	Content     string      `json:"content"`
	CreatedAt   time.Time   `json:"createdAt"`
	IsDeleted   bool        `json:"isDeleted"`
	Section     NewsSection `json:"section" gorm:"foreignKey:SectionID"`
	File        File        `json:"image" gorm:"foreignKey:FileID"`
}
