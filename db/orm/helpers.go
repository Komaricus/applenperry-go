package orm

import (
	"gorm.io/gorm"
)

type Filters struct {
	Search     string
	SortColumn string
	SortOrder  string
}

func GetList(tx *gorm.DB, target interface{}, filter Filters) error {
	q := tx.Model(target)

	if filter.Search != "" {
		filter.Search = "%" + filter.Search + "%"
		q.Where("name LIKE ?", filter.Search)
	}

	if filter.SortColumn == "created_at" && filter.SortOrder == "asc" {
		q.Order("created_at asc")
	}

	if filter.SortColumn == "created_at" && filter.SortOrder == "desc" {
		q.Order("created_at desc")
	}

	if filter.SortColumn == "updated_at" && filter.SortOrder == "asc" {
		q.Order("updated_at asc")
	}

	if filter.SortColumn == "updated_at" && filter.SortOrder == "desc" {
		q.Order("updated_at desc")
	}

	if err := q.Find(target).Error; err != nil {
		return err
	}

	return nil
}

func GetFirst(tx *gorm.DB, target interface{}, id string) error {
	if err := tx.Where("id = ?", id).First(target).Error; err != nil {
		return err
	}

	return nil
}
