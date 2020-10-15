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
	CreatedAt time.Time `json:"createdAt"`
	IsDeleted bool      `json:"isDeleted"`
}
