package model

import "time"

type FilesTable struct{}

func (FilesTable) TableName() string {
	return "dbo.files"
}

type File struct {
	FilesTable
	ID           string    `json:"id" gorm:"primarykey"`
	FileName     string    `json:"fileName"`
	Path         string    `json:"path"`
	OriginalName string    `json:"originalName"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type GetFilesResponse struct {
	Total int64  `json:"total"`
	Files []File `json:"files"`
}
