package api

import (
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

func GetNewProducts(c *gin.Context) {
	var products []model.ProductsListResponse
	if err := db.DB.Order("created_at desc").Limit(10).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, p := range products {
		var paf []model.ProductsAndFiles
		if err := db.DB.Where("product_id = ?", p.ID).Order("priority").Preload("File").Find(&paf).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(paf) > 0 {
			products[i].MainImage = paf[0].File
		}
	}

	c.JSON(http.StatusOK, products)
}

func GetProducts(c *gin.Context) {
	var products []model.Product
	q := db.DB.Preload("ProductsType").Preload("ProductsSugarType").Preload("Vendor").Preload("Vendor.File")
	if err := orm.GetList(q, &products, orm.Filters{
		Search:     c.Query("search"),
		SortColumn: c.Query("sort"),
		SortOrder:  c.Query("order"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i, p := range products {
		var paf []model.ProductsAndFiles
		if err := db.DB.Where("product_id = ?", p.ID).Order("priority").Preload("File").Find(&paf).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, f := range paf {
			products[i].Files = append(products[i].Files, f.File)
		}
		if len(products[i].Files) > 0 {
			products[i].MainImage = products[i].Files[0]
		}

		var pac []model.ProductsAndCategories
		if err := db.DB.Where("product_id = ?", p.ID).Order("priority").Preload("Category").Find(&pac).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, cat := range pac {
			products[i].Categories = append(products[i].Categories, cat.Category)
		}
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}
	var p model.Product
	q := db.DB.Preload("ProductsType").Preload("ProductsSugarType").Preload("Vendor").Preload("Vendor.File")
	if err := orm.GetFirst(q, &p, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var paf []model.ProductsAndFiles
	if err := db.DB.Where("product_id = ?", p.ID).Order("priority").Preload("File").Find(&paf).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, f := range paf {
		p.Files = append(p.Files, f.File)
	}
	if len(p.Files) > 0 {
		p.MainImage = p.Files[0]
	}

	var pac []model.ProductsAndCategories
	if err := db.DB.Where("product_id = ?", p.ID).Order("priority").Preload("Category").Find(&pac).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, cat := range pac {
		p.Categories = append(p.Categories, cat.Category)
	}

	c.JSON(http.StatusOK, p)
}

func CreateProduct(c *gin.Context) {
	var p model.Product
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	p.ID = id.String()

	if err := db.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updatePriorities(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func UpdateProduct(c *gin.Context) {
	var p model.Product
	if err := c.Bind(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("product_id = ?", p.ID).Delete(model.ProductsAndFiles{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("product_id = ?", p.ID).Delete(model.ProductsAndCategories{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Updates(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := updatePriorities(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	if err := db.DB.Where("product_id = ?", id).Delete(model.ProductsAndFiles{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("product_id = ?", id).Delete(model.ProductsAndCategories{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Delete(model.Product{ID: id}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deleted"})
}

func updatePriorities(product model.Product) error {
	for i, cat := range product.Categories {
		update := &model.ProductsAndCategories{
			ProductID:  product.ID,
			CategoryID: cat.ID,
			Priority:   i,
		}
		if err := db.DB.Updates(&update).Error; err != nil {
			return err
		}
	}

	for i, f := range product.Files {
		update := &model.ProductsAndFiles{
			ProductID: product.ID,
			FileID:    f.ID,
			Priority:  i,
		}
		if err := db.DB.Updates(&update).Error; err != nil {
			return err
		}
	}

	return nil
}
