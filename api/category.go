package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"net/http"
)

func GetCategoryByURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url param required"})
		return
	}
	var category model.Category
	if err := db.DB.Where("url = ?", url).First(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func GetCategoriesWithChild(c *gin.Context) {
	var categories []model.CategoryWithChild
	q := db.DB.Preload("Child", func(db *gorm.DB) *gorm.DB {
		return db.Order("name")
	}).Order("created_at desc").Where("parent_id is null")
	if err := q.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	countryCategoryID := "countries"
	countryCategory := model.CategoryWithChild{
		ID:       countryCategoryID,
		Name:     "Регион",
		URL:      countryCategoryID,
		ParentID: nil,
		Child:    nil,
	}

	var countries []model.Country
	if err := db.DB.Order("name").Find(&countries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, country := range countries {
		countryCategory.Child = append(countryCategory.Child, model.CategoryWithChild{
			ID:       country.ID,
			Name:     country.Name,
			URL:      country.URL,
			ParentID: &countryCategoryID,
			Child:    nil,
		})
	}

	categories = append(categories, countryCategory)

	productTypeID := "product-type"
	productTypeCategory := model.CategoryWithChild{
		ID:       productTypeID,
		Name:     "Тип",
		URL:      productTypeID,
		ParentID: nil,
		Child:    nil,
	}

	var productsTypes []model.ProductsType
	if err := db.DB.Order("name").Find(&productsTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, pt := range productsTypes {
		productTypeCategory.Child = append(productTypeCategory.Child, model.CategoryWithChild{
			ID:       pt.ID,
			Name:     pt.Name,
			URL:      pt.URL,
			ParentID: &productTypeID,
			Child:    nil,
		})
	}

	categories = append(categories, productTypeCategory)

	sugarTypeID := "sugar-type"
	sugarTypeCategory := model.CategoryWithChild{
		ID:       sugarTypeID,
		Name:     "Сахар",
		URL:      sugarTypeID,
		ParentID: nil,
		Child:    nil,
	}

	var sugarTypes []model.ProductsSugarType
	if err := db.DB.Order("name").Find(&sugarTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, st := range sugarTypes {
		sugarTypeCategory.Child = append(sugarTypeCategory.Child, model.CategoryWithChild{
			ID:       st.ID,
			Name:     st.Name,
			URL:      st.URL,
			ParentID: &sugarTypeID,
			Child:    nil,
		})
	}

	categories = append(categories, sugarTypeCategory)

	c.JSON(http.StatusOK, categories)
}

func GetCategories(c *gin.Context) {
	var categories []model.Category
	if err := orm.GetList(db.DB, &categories, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var category model.Category
	if err := orm.GetFirst(db.DB, &category, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	category.ID = id.String()

	if err := db.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	var category model.Category
	if err := c.Bind(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if category.ParentID == nil {
		if err := db.DB.Model(category).UpdateColumn("parent_id", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Delete(model.Category{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func GetPossibleToDeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var categories []model.Category
	if err := db.DB.Where("parent_id = ?", id).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var pac []model.ProductsAndCategories
	if err := db.DB.Where("category_id = ?", id).Find(&pac).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ids := make([]string, 0, len(pac))
	for _, n := range pac {
		ids = append(ids, n.ProductID)
	}

	var products []model.Product
	if err := db.DB.Where("id IN (?)", ids).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(categories) > 0 || len(products) > 0 {
		deleteConflicts := make(map[string]interface{})
		deleteConflicts["products"] = products
		deleteConflicts["categories"] = categories

		c.JSON(http.StatusOK, gin.H{
			"id":              id,
			"status":          "not_deletable",
			"deleteConflicts": deleteConflicts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deletable"})
}
