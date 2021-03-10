package entities

import "time"

// OrdersStatus int
type OrdersStatus int

// OrdersStatus Master
const (
	Pending OrdersStatus = iota + 1
	Approve
	Reject
	Cancel
)

// Orders struct
type Orders struct {
	ID          int          `db:"id" json:"id"`
	UserID      int          `db:"user_id" json:"user_id"`
	ProductID   int          `db:"product_id" json:"product_id"`
	ProductName string       `db:"product_name" json:"product_name"`
	Price       float32      `db:"price" json:"price"`
	Qty         int          `db:"qty" json:"qty"`
	TotalPrice  float32      `db:"total_price" json:"total_price"`
	Status      OrdersStatus `db:"status" json:"status"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}
