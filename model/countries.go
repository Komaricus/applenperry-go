package model

import "time"

type CountriesTable struct{}

func (CountriesTable) TableName() string {
	return "dbo.countries"
}

type Country struct {
	CountriesTable
	ID        string    `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	Flag      string    `json:"flag"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Icon      string    `json:"icon"`

	File     File `json:"image" gorm:"foreignKey:Flag"`
	IconFile File `json:"iconFile" gorm:"foreignKey:Icon"`
}
