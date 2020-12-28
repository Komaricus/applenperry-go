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
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Icon      string    `json:"icon"`

	IconFile File `json:"iconFile" gorm:"foreignKey:Icon"`
}
