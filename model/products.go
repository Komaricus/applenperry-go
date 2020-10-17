package model

import (
	"time"
)

type ProductsTable struct{}

func (ProductsTable) TableName() string {
	return "dbo.products"
}

type Product struct {
	ProductsTable
	ID          string  `json:"id" gorm:"primarykey"`
	Name        string  `json:"name"`
	Subheader   string  `json:"subheader"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	TypeID      string  `json:"typeId" gorm:"column:type"`
	Amount      uint    `json:"amount"`
	Size        float64 `json:"size"`
	Strength    float64 `json:"strength"`
	SugarTypeID string  `json:"sugarTypeId" gorm:"column:sugar_type"`
	Price       float64 `json:"price"`
	VendorCode  string  `json:"vendorCode"`
	VendorID    string  `json:"vendorId"`

	CreatedAt time.Time `json:"createdAt"`
	IsDeleted bool      `json:"isDeleted"`

	MainImage         File              `json:"image" gorm:"-"`
	ProductsType      ProductsType      `json:"productsType" gorm:"foreignKey:TypeID"`
	ProductsSugarType ProductsSugarType `json:"productsSugarType" gorm:"foreignKey:SugarTypeID"`
	Vendor            Vendor            `json:"vendor" gorm:"foreignKey:VendorID"`
	Files             []File            `json:"images" gorm:"many2many:dbo.products_and_files;joinForeignKey:ProductID;JoinReferences:FileID"`
	Categories        []Category        `json:"categories" gorm:"many2many:dbo.products_and_categories;joinForeignKey:ProductID;JoinReferences:CategoryID"`
}

type ProductsAndFilesTable struct{}

func (ProductsAndFilesTable) TableName() string {
	return "dbo.products_and_files"
}

type ProductsAndFiles struct {
	ProductsAndFilesTable
	ProductID string `gorm:"primaryKey"`
	FileID    string `gorm:"primaryKey"`
	Priority  int
	File      File `gorm:"foreignKey:FileID"`
}

type ProductsAndCategoriesTable struct{}

func (ProductsAndCategoriesTable) TableName() string {
	return "dbo.products_and_categories"
}

type ProductsAndCategories struct {
	ProductsAndCategoriesTable
	ProductID  string `gorm:"primaryKey"`
	CategoryID string `gorm:"primaryKey"`
	Priority   int
	Category   Category `gorm:"foreignKey:CategoryID"`
}
