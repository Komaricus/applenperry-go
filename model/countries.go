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
	CreatedAt time.Time `json:"createdAt"`
	IsDeleted bool      `json:"isDeleted"`
	File      File      `json:"image" gorm:"foreignKey:Flag"`
}
