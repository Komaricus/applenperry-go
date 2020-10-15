package model

import "time"

type ProductsSugarTypesTable struct{}

func (ProductsSugarTypesTable) TableName() string {
	return "dbo.products_sugar_types"
}

type ProductsSugarType struct {
	ProductsSugarTypesTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	IsDeleted bool      `json:"isDeleted"`
}
