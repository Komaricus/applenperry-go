package model

import "time"

type CategoriesTable struct{}

func (CategoriesTable) TableName() string {
	return "dbo.categories"
}

type Category struct {
	CategoriesTable
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	ParentID    *string   `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	IsDeleted   bool      `json:"isDeleted"`
}
