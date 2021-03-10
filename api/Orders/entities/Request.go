package entities

// OrdersListUsersRequest struct
type OrdersListUsersRequest struct {
	UserID int    `json:"user_id" validate:"required"`
	Limit  string `json:"limit" validate:"required"`
	Offset string `json:"offset" validate:"required"`
}

// OrdersCreateRequest struct
type OrdersCreateRequest struct {
	UserID      int     `json:"user_id" validate:"required"`
	ProductID   int     `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
	Qty         int     `json:"qty" validate:"required"`
}

// OrdersUpdateRequest struct
type OrdersUpdateRequest struct {
	UserID  int `json:"user_id" validate:"required"`
	OrderID int `json:"order_id" validate:"required"`
	Qty     int `json:"qty" validate:"required"`
}

// OrdersCancelRequest struct
type OrdersCancelRequest struct {
	UserID  int `json:"user_id" validate:"required"`
	OrderID int `json:"order_id" validate:"required"`
}

// OrdersListAdminRequest struct
type OrdersListAdminRequest struct {
	Limit     string `json:"limit" validate:"required"`
	Offset    string `json:"offset" validate:"required"`
	Status    int    `json:"status" validate:"-"`
	ProductID int    `json:"product_id" validate:"-"`
}

// OrdersApproveRequest struct
type OrdersApproveRequest struct {
	UserID  int `json:"user_id" validate:"required"`
	OrderID int `json:"order_id" validate:"required"`
}

// OrdersRejectRequest struct
type OrdersRejectRequest struct {
	UserID  int `json:"user_id" validate:"required"`
	OrderID int `json:"order_id" validate:"required"`
}

// GetProductsByIDPayload struct
type GetProductsByIDPayload struct {
	ProductID int `json:"product_id"`
}

// GetProductsByIDResponse struct
type GetProductsByIDResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    *Products `json:"data"`
}

// ProductsUpdatePayload struct
type ProductsUpdatePayload struct {
	UserID    int `json:"user_id" validate:"required"`
	ProductID int `json:"product_id" validate:"required"`
	Qty       int `json:"qty,omitempty" validate:"-"`
}

// ProductsUpdateResponse struct
type ProductsUpdateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
