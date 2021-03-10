package usecases

import (
	"context"
	"strconv"
	"time"

	"github.com/mrdhira/warpin-test/api/Orders/entities"
	"github.com/mrdhira/warpin-test/api/Orders/infrastructures/database"
	"github.com/mrdhira/warpin-test/api/Orders/infrastructures/repositories"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

// IOrdersUsecases interface
type IOrdersUsecases interface {
	OrdersListUsers(ctx context.Context, Data *entities.OrdersListUsersRequest) (Response *pkg.JSONResponse, err error)
	OrdersCreate(ctx context.Context, Data *entities.OrdersCreateRequest) (Response *pkg.JSONResponse, err error)
	OrdersUpdate(ctx context.Context, Data *entities.OrdersUpdateRequest) (Response *pkg.JSONResponse, err error)
	OrdersCancel(ctx context.Context, Data *entities.OrdersCancelRequest) (Response *pkg.JSONResponse, err error)
	OrdersListAdmin(ctx context.Context, Data *entities.OrdersListAdminRequest) (Response *pkg.JSONResponse, err error)
	OrdersApprove(ctx context.Context, Data *entities.OrdersApproveRequest) (Response *pkg.JSONResponse, err error)
	OrdersReject(ctx context.Context, Data *entities.OrdersRejectRequest) (Response *pkg.JSONResponse, err error)
}

// OrdersUsecases struct
type OrdersUsecases struct {
	OrdersRepository   repositories.IOrdersRepository
	ProductsRepository repositories.IProductsrepository
}

// InitOrdersUsecases func
func InitOrdersUsecases() *OrdersUsecases {
	// Init Repositories
	ordersRepository := new(repositories.OrdersRepository)
	ordersRepository.PG = &database.PostgresConnection{}
	ordersRepository.Redis = &database.RedisConnection{}

	productsRepository := new(repositories.ProductsRepository)

	return &OrdersUsecases{
		OrdersRepository:   ordersRepository,
		ProductsRepository: productsRepository,
	}
}

// OrdersListUsers func
func (u *OrdersUsecases) OrdersListUsers(ctx context.Context, Data *entities.OrdersListUsersRequest) (Response *pkg.JSONResponse, err error) {
	Limit, err := strconv.Atoi(Data.Limit)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when parse to int for limit",
		}).Error(err)
		return
	}

	Offset, err := strconv.Atoi(Data.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when parse to int for offset",
		}).Error(err)
		return
	}

	Orders, err := u.OrdersRepository.OrdersFindByUserID(ctx, Limit, Offset, Data.UserID)
	if err != nil {
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "OK",
		Data:    Orders,
	}, nil
}

// OrdersCreate func
func (u *OrdersUsecases) OrdersCreate(ctx context.Context, Data *entities.OrdersCreateRequest) (Response *pkg.JSONResponse, err error) {
	GetProductsByIDPayload := &entities.GetProductsByIDPayload{
		ProductID: Data.ProductID,
	}
	Products, err := u.ProductsRepository.GetProductsByID(ctx, GetProductsByIDPayload)
	if err != nil {
		return
	}

	if Products.Status != entities.Active {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Product sedang tidak aktif, silahkan hubungi cs",
		}, nil
	}

	if Data.Qty > Products.Qty {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Kuantitas yang di order lebih banyak daripada stok yang tersedia",
		}, nil
	}

	Tx, err := u.OrdersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	Orders := &entities.Orders{
		UserID:      Data.UserID,
		ProductID:   Data.ProductID,
		ProductName: Data.ProductName,
		Price:       Data.Price,
		Qty:         Data.Qty,
		TotalPrice:  Data.Price * float32(Data.Qty),
		Status:      entities.Pending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	Orders.ID, err = u.OrdersRepository.OrdersStore(ctx, Tx, Orders)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	OrdersLog := &entities.OrdersLog{
		OrderID:     Orders.ID,
		UserID:      Orders.UserID,
		ProductID:   Orders.ProductID,
		ProductName: Orders.ProductName,
		Price:       Orders.Price,
		Qty:         Orders.Qty,
		TotalPrice:  Orders.TotalPrice,
		Status:      entities.Pending,
		Event:       entities.EventCreate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	OrdersLog.ID, err = u.OrdersRepository.OrdersLogStore(ctx, Tx, OrdersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	ProductsUpdatePayload := &entities.ProductsUpdatePayload{
		UserID:    0,
		ProductID: Data.ProductID,
		Qty:       Products.Qty - Data.Qty,
	}
	err = u.ProductsRepository.ProductsUpdate(ctx, ProductsUpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Orders berhasil dibuat",
	}, nil
}

// OrdersUpdate func
func (u *OrdersUsecases) OrdersUpdate(ctx context.Context, Data *entities.OrdersUpdateRequest) (Response *pkg.JSONResponse, err error) {
	Orders, err := u.OrdersRepository.OrdersFindByID(ctx, Data.OrderID)
	if err != nil {
		return
	}

	if Orders.Status != entities.Pending {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Order tidak sedang pending",
		}, nil
	}

	if Orders.UserID != Data.UserID {
		return &pkg.JSONResponse{
			Code:    403,
			Message: "Order bukan milik users",
		}, nil
	}

	Tx, err := u.OrdersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	OrdersLog := &entities.OrdersLog{
		OrderID:     Orders.ID,
		UserID:      Orders.UserID,
		ProductID:   Orders.ProductID,
		ProductName: Orders.ProductName,
		Price:       Orders.Price,
		Qty:         Data.Qty,
		TotalPrice:  Orders.Price * float32(Data.Qty),
		Status:      entities.Pending,
		Event:       entities.EventUpdate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{
		"qty":         Data.Qty,
		"total_price": Orders.Price * float32(Data.Qty),
		"updated_at":  time.Now(),
	}

	err = u.OrdersRepository.OrdersUpdate(ctx, Tx, Data.OrderID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	OrdersLog.ID, err = u.OrdersRepository.OrdersLogStore(ctx, Tx, OrdersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	GetProductsByIDPayload := &entities.GetProductsByIDPayload{
		ProductID: Orders.ProductID,
	}
	Products, err := u.ProductsRepository.GetProductsByID(ctx, GetProductsByIDPayload)
	if err != nil {
		return
	}

	ProductsUpdatePayload := &entities.ProductsUpdatePayload{
		UserID:    0,
		ProductID: Orders.ProductID,
		// Qty:       Products.Qty + Orders.Qty,
	}

	if Data.Qty >= Orders.Qty {
		ProductsUpdatePayload.Qty = Products.Qty - (Data.Qty - Orders.Qty)
	} else {
		ProductsUpdatePayload.Qty = Products.Qty + (Orders.Qty - Data.Qty)
	}
	err = u.ProductsRepository.ProductsUpdate(ctx, ProductsUpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Orders berhasil di update",
	}, nil
}

// OrdersCancel func
func (u *OrdersUsecases) OrdersCancel(ctx context.Context, Data *entities.OrdersCancelRequest) (Response *pkg.JSONResponse, err error) {
	Orders, err := u.OrdersRepository.OrdersFindByID(ctx, Data.OrderID)
	if err != nil {
		return
	}

	if Orders.Status != entities.Pending {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Order tidak sedang pending",
		}, nil
	}

	if Orders.UserID != Data.UserID {
		return &pkg.JSONResponse{
			Code:    403,
			Message: "Order bukan milik users",
		}, nil
	}

	Tx, err := u.OrdersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	OrdersLog := &entities.OrdersLog{
		OrderID:     Orders.ID,
		UserID:      Orders.UserID,
		ProductID:   Orders.ProductID,
		ProductName: Orders.ProductName,
		Price:       Orders.Price,
		Qty:         Orders.Qty,
		TotalPrice:  Orders.TotalPrice,
		Status:      entities.Cancel,
		Event:       entities.EventCancel,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{
		"status":     entities.Cancel,
		"updated_at": time.Now(),
	}

	err = u.OrdersRepository.OrdersUpdate(ctx, Tx, Data.OrderID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	OrdersLog.ID, err = u.OrdersRepository.OrdersLogStore(ctx, Tx, OrdersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	GetProductsByIDPayload := &entities.GetProductsByIDPayload{
		ProductID: Orders.ProductID,
	}
	Products, err := u.ProductsRepository.GetProductsByID(ctx, GetProductsByIDPayload)
	if err != nil {
		return
	}

	ProductsUpdatePayload := &entities.ProductsUpdatePayload{
		UserID:    0,
		ProductID: Orders.ProductID,
		Qty:       Products.Qty + Orders.Qty,
	}
	err = u.ProductsRepository.ProductsUpdate(ctx, ProductsUpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Orders berhasil di cancel",
	}, nil
}

// OrdersListAdmin func
func (u *OrdersUsecases) OrdersListAdmin(ctx context.Context, Data *entities.OrdersListAdminRequest) (Response *pkg.JSONResponse, err error) {
	Condition := map[string]interface{}{}
	if Data.Status != 0 {
		Condition["status"] = Data.Status
	}

	if Data.ProductID != 0 {
		Condition["product_id"] = Data.ProductID
	}

	Limit, err := strconv.Atoi(Data.Limit)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when parse to int for limit",
		}).Error(err)
		return
	}
	Offset, err := strconv.Atoi(Data.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when parse to int for offset",
		}).Error(err)
		return
	}

	Orders, err := u.OrdersRepository.OrdersFind(ctx, Limit, Offset, Condition)
	if err != nil {
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "OK",
		Data:    Orders,
	}, nil
}

// OrdersApprove func
func (u *OrdersUsecases) OrdersApprove(ctx context.Context, Data *entities.OrdersApproveRequest) (Response *pkg.JSONResponse, err error) {
	Orders, err := u.OrdersRepository.OrdersFindByID(ctx, Data.OrderID)
	if err != nil {
		return
	}

	if Orders.Status != entities.Pending {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Order tidak sedang pending",
		}, nil
	}

	Tx, err := u.OrdersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	OrdersLog := &entities.OrdersLog{
		OrderID:     Orders.ID,
		UserID:      Orders.UserID,
		ProductID:   Orders.ProductID,
		ProductName: Orders.ProductName,
		Price:       Orders.Price,
		Qty:         Orders.Qty,
		TotalPrice:  Orders.TotalPrice,
		Status:      entities.Approve,
		Event:       entities.EventApprove,
		AdminID:     Data.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{
		"status":     entities.Approve,
		"updated_at": time.Now(),
	}

	err = u.OrdersRepository.OrdersUpdate(ctx, Tx, Data.OrderID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	OrdersLog.ID, err = u.OrdersRepository.OrdersLogStore(ctx, Tx, OrdersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Orders berhasil di approve",
	}, nil
}

// OrdersReject func
func (u *OrdersUsecases) OrdersReject(ctx context.Context, Data *entities.OrdersRejectRequest) (Response *pkg.JSONResponse, err error) {
	Orders, err := u.OrdersRepository.OrdersFindByID(ctx, Data.OrderID)
	if err != nil {
		return
	}

	if Orders.Status != entities.Pending {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Order tidak sedang pending",
		}, nil
	}

	Tx, err := u.OrdersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	OrdersLog := &entities.OrdersLog{
		OrderID:     Orders.ID,
		UserID:      Orders.UserID,
		ProductID:   Orders.ProductID,
		ProductName: Orders.ProductName,
		Price:       Orders.Price,
		Qty:         Orders.Qty,
		TotalPrice:  Orders.TotalPrice,
		Status:      entities.Reject,
		Event:       entities.EventReject,
		AdminID:     Data.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{
		"status":     entities.Reject,
		"updated_at": time.Now(),
	}

	err = u.OrdersRepository.OrdersUpdate(ctx, Tx, Data.OrderID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	OrdersLog.ID, err = u.OrdersRepository.OrdersLogStore(ctx, Tx, OrdersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	GetProductsByIDPayload := &entities.GetProductsByIDPayload{
		ProductID: Orders.ProductID,
	}
	Products, err := u.ProductsRepository.GetProductsByID(ctx, GetProductsByIDPayload)
	if err != nil {
		return
	}

	ProductsUpdatePayload := &entities.ProductsUpdatePayload{
		UserID:    0,
		ProductID: Orders.ProductID,
		Qty:       Products.Qty + Orders.Qty,
	}
	err = u.ProductsRepository.ProductsUpdate(ctx, ProductsUpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Orders berhasil di reject",
	}, nil
}
