package entities

import "time"

// ProductsStatus int
type ProductsStatus int

// ProductsStatus Master
const (
	Active ProductsStatus = iota + 1
	InActive
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
