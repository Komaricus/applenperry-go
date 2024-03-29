package model

import "time"

type OrdersTable struct{}

func (OrdersTable) TableName() string {
	return "dbo.orders"
}

type GetOrder struct {
	OrdersTable
	ID        string    `json:"id" gorm:"primarykey"`
	Code      int       `json:"code"`
	UserName  string    `json:"userName"`
	UserPhone string    `json:"userPhone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Products []OrderAndProduct `json:"products" gorm:"foreignKey:OrderID"`
}

type CreateOrder struct {
	OrdersTable
	ID        string    `json:"id" gorm:"primarykey"`
	Code      int       `json:"code" gorm:"-"`
	UserName  string    `json:"userName"`
	UserPhone string    `json:"userPhone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Products []OrderAndProduct `json:"products" gorm:"-"`
}

type DeleteOrder struct {
	OrdersTable
	ID string `json:"id" gorm:"primarykey"`
}

type UpdateOrder struct {
	OrdersTable
	ID        string    `json:"id" gorm:"primarykey"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type OrdersAndProductsTable struct{}

func (OrdersAndProductsTable) TableName() string {
	return "dbo.orders_and_products"
}

type OrderAndProduct struct {
	OrdersAndProductsTable
	OrderID      string `json:"orderId"`
	ProductID    string `json:"productId"`
	ProductCount uint   `json:"productCount"`

	Product Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Order   GetOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
