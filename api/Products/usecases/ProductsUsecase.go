package usecases

import (
	"context"
	"strconv"
	"time"

	"github.com/mrdhira/warpin-test/api/Products/entities"
	"github.com/mrdhira/warpin-test/api/Products/infrastructures/database"
	"github.com/mrdhira/warpin-test/api/Products/infrastructures/repositories"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
)

// IProductsUsecases interface
type IProductsUsecases interface {
	GetProducts(ctx context.Context, Data *entities.GetProductsRequest) (Response *pkg.JSONResponse, err error)
	GetProductsByID(ctx context.Context, Data *entities.GetProductsByIDRequest) (Response *pkg.JSONResponse, err error)
	AddProducts(ctx context.Context, Data *entities.AddProductsRequest) (Response *pkg.JSONResponse, err error)
	UpdateProducts(ctx context.Context, Data *entities.UpdateProductsRequest) (Response *pkg.JSONResponse, err error)
}

// ProductsUsecases struct
type ProductsUsecases struct {
	ProductsRepository repositories.IProductsRepository
	OrdersRepository   repositories.IOrdersRepository
}

// InitProductsUsecases func
func InitProductsUsecases() *ProductsUsecases {
	// Init Repositories
	productsRepository := new(repositories.ProductsRepository)
	productsRepository.PG = &database.PostgresConnection{}
	productsRepository.Redis = &database.RedisConnection{}

	ordersRepository := new(repositories.OrdersRepository)

	return &ProductsUsecases{
		ProductsRepository: productsRepository,
		OrdersRepository:   ordersRepository,
	}
}

// GetProducts usecases
func (u *ProductsUsecases) GetProducts(ctx context.Context, Data *entities.GetProductsRequest) (Response *pkg.JSONResponse, err error) {
	Condition := map[string]interface{}{}
	if Data.Status != 0 {
		Condition["status"] = Data.Status
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

	Products, err := u.ProductsRepository.ProductsFind(ctx, Limit, Offset, Condition)
	if err != nil {
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "OK",
		Data:    Products,
	}, nil
}

// GetProductsByID usecases
func (u *ProductsUsecases) GetProductsByID(ctx context.Context, Data *entities.GetProductsByIDRequest) (Response *pkg.JSONResponse, err error) {
	Products, err := u.ProductsRepository.ProductsFindOneByID(ctx, Data.ProductID)
	if err != nil {
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "OK",
		Data:    Products,
	}, nil
}

// AddProducts usecases
func (u *ProductsUsecases) AddProducts(ctx context.Context, Data *entities.AddProductsRequest) (Response *pkg.JSONResponse, err error) {
	Tx, err := u.ProductsRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	Products := &entities.Products{
		Name:      Data.Name,
		Price:     Data.Price,
		Qty:       Data.Qty,
		Status:    entities.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	Products.ID, err = u.ProductsRepository.ProductsStore(ctx, Tx, Products)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	ProductsLog := &entities.ProductsLog{
		ProductID: Products.ID,
		UserID:    Data.UserID,
		Name:      Data.Name,
		Price:     Data.Price,
		Qty:       Data.Qty,
		Status:    entities.Active,
		Event:     entities.EventInsert,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ProductsLog.ID, err = u.ProductsRepository.ProductsLogStore(ctx, Tx, ProductsLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Produk berhasil ditambahkan",
		Data:    Products,
	}, nil
}

// UpdateProducts usecases
func (u *ProductsUsecases) UpdateProducts(ctx context.Context, Data *entities.UpdateProductsRequest) (Response *pkg.JSONResponse, err error) {
	if Data.Name == "" && Data.Price == nil && Data.Qty == nil && Data.Status == nil {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Nama, harga, dan kuantitas tidak bisa kosong semua",
		}, nil
	}

	Products, err := u.ProductsRepository.ProductsFindOneByID(ctx, Data.ProductID)
	if err != nil {
		return
	}

	if Products == nil {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Product dengan ID " + strconv.Itoa(Data.ProductID) + "tidak ditemukan",
		}, nil
	}

	if Data.Status != nil {
		if entities.ProductsStatus(*Data.Status) == entities.InActive {
			GetOrdersByPrductIDPayload := &entities.GetOrdersByPrductIDPayload{
				Limit:     1,
				Offset:    0,
				Status:    1, // Pending Status
				ProductID: Data.ProductID,
			}
			Orders, err := u.OrdersRepository.GetOrdersByProductID(ctx, GetOrdersByPrductIDPayload)
			if err != nil {
				return nil, err
			}

			if len(Orders) != 0 {
				return &pkg.JSONResponse{
					Code:    422,
					Message: "Masih terdapat orders yang pending, tidak bisa menonaktifkan produk, silahkan mengurangi qty produk terlebih dahulu atau ubah status order",
				}, nil
			}
		}
	}

	Tx, err := u.ProductsRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	ProductsLog := &entities.ProductsLog{
		ProductID: Products.ID,
		UserID:    Data.UserID,
		Name:      Products.Name,
		Price:     Products.Price,
		Qty:       Products.Qty,
		Status:    Products.Status,
		Event:     entities.EventUpdate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	UpdatePayload := map[string]interface{}{}
	if Data.Name != "" {
		ProductsLog.Name = Data.Name
		UpdatePayload["name"] = Data.Name
	}
	if Data.Price != nil {
		ProductsLog.Price = *Data.Price
		UpdatePayload["price"] = Data.Price
	}
	if Data.Qty != nil {
		ProductsLog.Qty = *Data.Qty
		UpdatePayload["qty"] = Data.Qty
	}
	if Data.Status != nil {
		ProductsLog.Status = entities.ProductsStatus(*Data.Status)
		UpdatePayload["status"] = entities.ProductsStatus(*Data.Status)
	}

	err = u.ProductsRepository.ProductsUpdate(ctx, Tx, Data.ProductID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	ProductsLog.ID, err = u.ProductsRepository.ProductsLogStore(ctx, Tx, ProductsLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Products berhasil diupdate",
	}, nil
}
