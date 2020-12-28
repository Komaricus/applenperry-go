package model

import "time"

type CategoriesTable struct{}

func (CategoriesTable) TableName() string {
	return "dbo.categories"
}

type Category struct {
	CategoriesTable
	ID          string    `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	ParentID    *string   `json:"parentId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Icon        string    `json:"icon"`

	IconFile File `json:"iconFile" gorm:"foreignKey:Icon"`
}

type CategoryWithChild struct {
	CategoriesTable
	ID       string  `json:"id" gorm:"primarykey"`
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	ParentID *string `json:"parentId"`
	Icon     string  `json:"icon"`

	IconFile File                `json:"iconFile" gorm:"foreignKey:Icon"`
	Child    []CategoryWithChild `json:"child" gorm:"foreignKey:ParentID"`
}
