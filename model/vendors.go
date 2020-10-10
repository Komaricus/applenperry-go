package model

import "time"

type VendorsTable struct{}

func (VendorsTable) TableName() string {
	return "dbo.vendors"
}

type Vendor struct {
	VendorsTable
	ID          string    `json:"id" gorm:"primarykey"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	FileID      string    `json:"fileId"`
	Description string    `json:"description"`
	CountryID   string    `json:"countryId"`
	CreatedAt   time.Time `json:"createdAt"`
	IsDeleted   bool      `json:"isDeleted"`
	File        File      `json:"image" gorm:"foreignKey:FileID"`
	Country     Country   `json:"country" gorm:"foreignKey:CountryID"`
}
