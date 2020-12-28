package api

import (
	"errors"
	"github.com/applenperry-go/db"
	"github.com/applenperry-go/db/orm"
	"github.com/applenperry-go/model"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strconv"
)

func GetProductByURL(c *gin.Context) {
	url := c.Param("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url param required"})
		return
	}
	var p model.Product
	q := db.DB.Preload("ProductsType").Preload("ProductsSugarType").Preload("Vendor").Preload("Vendor.Country").Preload("Vendor.Country.File").Preload("MainImage")
	if err := q.Where("url = ?", url).First(&p).Error; err != nil {
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

func GetProductsWithPaginate(c *gin.Context) {
	// pagination
	page := 1
	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	perPage := -1
	perPageParam := c.Query("perPage")
	if perPageParam != "" {
		perPage, _ = strconv.Atoi(perPageParam)
	}
	offset := (page - 1) * perPage

	q := db.DB.Offset(offset)
	t := db.DB.Model(model.Product{}).Group("dbo.products.id")

	if perPage != -1 {
		q.Limit(perPage)
	}

	// search
	if search := c.Query("search"); search != "" {
		search = "%" + search + "%"
		q.Where("name LIKE ?", search)
		t.Where("name LIKE ?", search)
	}

	// new product flag
	if newProduct := c.Query("newProduct"); newProduct == "true" {
		q.Where("new_product is true")
		t.Where("new_product is true")
	}

	// sort
	sort := c.Query("sort")
	column := c.Query("column")
	if column == "created_at" && sort == "desc" {
		q.Order("created_at desc")
	}
	if column == "created_at" && sort == "asc" {
		q.Order("created_at asc")
	}
	if column == "price" && sort == "desc" {
		q.Order("price desc")
	}
	if column == "price" && sort == "asc" {
		q.Order("price asc")
	}

	q.Order("dbo.products.name")

	// filters
	if categoryUrl := c.Query("category"); categoryUrl != "" {
		var category model.Category
		if err := db.DB.Where("url = ?", categoryUrl).First(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var pac []model.ProductsAndCategories
		if err := db.DB.Where("category_id = ?", category.ID).Find(&pac).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		productIDs := make([]string, 0, len(pac))
		for _, p := range pac {
			productIDs = append(productIDs, p.ProductID)
		}

		q.Where("id IN (?)", productIDs)
		t.Where("id IN (?)", productIDs)
	}

	if countryUrl := c.Query("country"); countryUrl != "" {
		var country model.Country
		if err := db.DB.Where("url = ?", countryUrl).First(&country).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		q.Joins("Vendor").Where("\"Vendor\".country_id = ?", country.ID)
		t.Joins("Vendor").Where("\"Vendor\".country_id = ?", country.ID)
	}

	if typeUrl := c.Query("type"); typeUrl != "" {
		q.Joins("ProductsType").Where("\"ProductsType\".url = ?", typeUrl)
		t.Joins("ProductsType").Where("\"ProductsType\".url = ?", typeUrl)
	}

	if sugarUrl := c.Query("sugar"); sugarUrl != "" {
		q.Joins("ProductsSugarType").Where("\"ProductsSugarType\".url = ?", sugarUrl)
		t.Joins("ProductsSugarType").Where("\"ProductsSugarType\".url = ?", sugarUrl)
	}

	if vendorUrl := c.Query("vendor"); vendorUrl != "" {
		q.Joins("Vendor").Where("\"Vendor\".url = ?", vendorUrl)
		t.Joins("Vendor").Where("\"Vendor\".url = ?", vendorUrl)
	}

	if except := c.Query("except"); except != "" {
		q.Where("dbo.products.id NOT IN (?)", except)
		t.Where("dbo.products.id NOT IN (?)", except)
	}

	var products []model.Product
	if err := q.Preload("MainImage").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := t.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GetProductsResponse{
		Products: products,
		Total:    count,
	})
}

func GetProducts(c *gin.Context) {
	var products []model.Product
	q := db.DB.Preload("ProductsType").Preload("ProductsSugarType").Preload("Vendor").Preload("Vendor.File").Preload("MainImage")
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
	q := db.DB.Preload("ProductsType").Preload("ProductsSugarType").Preload("Vendor").Preload("Vendor.File").Preload("MainImage")
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

	if len(p.Files) == 0 || len(p.Categories) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("product images and categories can not be empty")})
		return
	}
	p.FileID = p.Files[0].ID

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

	if len(p.Files) == 0 || len(p.Categories) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("product images and categories can not be empty")})
		return
	}
	p.FileID = p.Files[0].ID

	if err := db.DB.Updates(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if p.Amount == 0 {
		if err := db.DB.Model(model.Product{}).Where("id = ?", p.ID).Update("amount", 0).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

func GetPossibleToDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param required"})
		return
	}

	var orders []model.OrderAndProduct
	if err := db.DB.Where("product_id = ?", id).Preload("Order").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(orders) > 0 {
		deleteConflicts := make(map[string]interface{})
		deleteConflicts["orders"] = orders

		c.JSON(http.StatusOK, gin.H{
			"id":              id,
			"status":          "not_deletable",
			"deleteConflicts": deleteConflicts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": "deletable"})
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
