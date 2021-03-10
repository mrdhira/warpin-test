package entities

import (
	"time"
)

// ProductsStatus int
type ProductsStatus int

// ProductsStatus Master
const (
	Active ProductsStatus = iota + 1
	InActive
)

// ProductsEvent string
type ProductsEvent string

// ProductsEvent Master
const (
	EventInsert ProductsEvent = "INSERT"
	EventUpdate ProductsEvent = "UPDATE"
)

// Products struct
type Products struct {
	ID        int            `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Price     float32        `db:"price" json:"price"`
	Qty       int            `db:"qty" json:"qty"`
	Status    ProductsStatus `db:"status" json:"status"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

// ProductsLog struct
type ProductsLog struct {
	ID        int            `db:"id" json:"id"`
	ProductID int            `db:"product_id" json:"product_id"`
	UserID    int            `db:"user_id" json:"user_id"`
	Name      string         `db:"name" json:"name"`
	Price     float32        `db:"price" json:"price"`
	Qty       int            `db:"qty" json:"qty"`
	Status    ProductsStatus `db:"status" json:"status"`
	Event     ProductsEvent  `db:"event" json:"event"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}
