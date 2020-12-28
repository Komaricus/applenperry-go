package model

import "time"

type ProductsTypesTable struct{}

func (ProductsTypesTable) TableName() string {
	return "dbo.products_types"
}

type ProductsType struct {
	ProductsTypesTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Icon      string    `json:"icon"`

	IconFile File `json:"iconFile" gorm:"foreignKey:Icon"`
}
