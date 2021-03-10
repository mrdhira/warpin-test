package entities

import (
	"time"
)

// OrdersStatus int
type OrdersStatus int

// OrdersStatus Master
const (
	Pending OrdersStatus = iota + 1
	Approve
	Reject
	Cancel
)

// OrdersEvent string
type OrdersEvent string

// OrdersEvent Master
const (
	EventCreate  OrdersEvent = "CREATE"
	EventApprove OrdersEvent = "APPROVE"
	EventReject  OrdersEvent = "REJECT"
	EventCancel  OrdersEvent = "CANCEL"
	EventUpdate  OrdersEvent = "UPDATE"
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

// OrdersLog struct
type OrdersLog struct {
	ID          int          `db:"id" json:"id"`
	OrderID     int          `db:"order_id" json:"order_id"`
	UserID      int          `db:"user_id" json:"user_id"`
	ProductID   int          `db:"product_id" json:"product_id"`
	ProductName string       `db:"product_name" json:"product_name"`
	Price       float32      `db:"price" json:"price"`
	Qty         int          `db:"qty" json:"qty"`
	TotalPrice  float32      `db:"total_price" json:"total_price"`
	Status      OrdersStatus `db:"status" json:"status"`
	Event       OrdersEvent  `db:"event" json:"event"`
	AdminID     int          `db:"admin_id" json:"admin_id"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}

// OrdersItems struct
// type OrdersItems struct {
// 	ID          int               `db:"id" json:"id"`
// 	OrderID     int               `db:"order_id" json:"order_id"`
// 	ProductID   int               `db:"product_id" json:"product_id"`
// 	ProductName string            `db:"product_name" json:"product_name"`
// 	Price       float32           `db:"price" json:"price"`
// 	Qty         int               `db:"qty" json:"qty"`
// 	TotalPrice  float32           `db:"total_price" json:"total_price"`
// 	Status      OrdersItemsStatus `db:"status" json:"status"`
// 	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
// 	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
// }

// OrdersItemsLog struct
// type OrdersItemsLog struct {
// 	ID          int               `db:"id" json:"id"`
// 	OrderItemID int               `db:"order_item_id" json:"order_item_id"`
// 	OrderID     int               `db:"order_id" json:"order_id"`
// 	ProductID   int               `db:"product_id" json:"product_id"`
// 	ProductName string            `db:"product_name" json:"product_name"`
// 	Price       float32           `db:"price" json:"price"`
// 	Qty         int               `db:"qty" json:"qty"`
// 	TotalPrice  float32           `db:"total_price" json:"total_price"`
// 	Status      OrdersItemsStatus `db:"status" json:"status"`
// 	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
// 	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
// }
