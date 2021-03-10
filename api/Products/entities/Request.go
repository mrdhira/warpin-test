package entities

// GetProductsRequest struct
type GetProductsRequest struct {
	Limit  string `json:"limit" validate:"required"`
	Offset string `json:"offset" validate:"required"`
	Status int    `json:"status" validate:"-"`
}

// GetProductsByIDRequest struct
type GetProductsByIDRequest struct {
	ProductID int `json:"product_id" validate:"required"`
}

// AddProductsRequest struct
type AddProductsRequest struct {
	UserID int     `json:"user_id" validate:"required"`
	Name   string  `json:"name" validate:"required"`
	Price  float32 `json:"price" validate:"required"`
	Qty    int     `json:"qty" validate:"required"`
}

// UpdateProductsRequest struct
type UpdateProductsRequest struct {
	UserID    int      `json:"user_id" validate:"-"`
	ProductID int      `json:"product_id" validate:"required"`
	Name      string   `json:"name" validate:"-"`
	Price     *float32 `json:"price,omitempty" validate:"-"`
	Qty       *int     `json:"qty,omitempty" validate:"-"`
	Status    *int     `json:"status,omitempty" validate:"-"`
}

// GetOrdersByPrductIDPayload struct
type GetOrdersByPrductIDPayload struct {
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
	Status    int `json:"status"`
	ProductID int `json:"product_id"`
}

// GetOrdersByPrductIDResponse struct
type GetOrdersByPrductIDResponse struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    []*Orders `json:"data"`
}
