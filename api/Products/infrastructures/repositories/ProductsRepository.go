package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
	dbr "github.com/gocraft/dbr/v2"
	"github.com/mrdhira/warpin-test/api/Products/entities"
	"github.com/mrdhira/warpin-test/api/Products/infrastructures/database"
	log "github.com/sirupsen/logrus"
)

// IProductsRepository interface
type IProductsRepository interface {
	Tx() (tx *dbr.Tx, err error)
	ProductsFind(ctx context.Context, Limit int, Offset int, Condition map[string]interface{}) (Products []*entities.Products, err error)
	ProductsFindOneByID(ctx context.Context, ID int) (Products *entities.Products, err error)
	ProductsStore(ctx context.Context, db *dbr.Tx, Products *entities.Products) (ID int, err error)
	ProductsLogStore(ctx context.Context, db *dbr.Tx, ProductsLog *entities.ProductsLog) (ID int, err error)
	ProductsUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error)
}

// ProductsRepository struct
type ProductsRepository struct {
	PG    database.IPostgresConnection
	Redis database.IRedisConnection
}

// Tx func to create new transaction
func (r *ProductsRepository) Tx() (tx *dbr.Tx, err error) {
	db := r.PG.PostgresTrade()

	tx, err = db.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when begin transaction in postgres",
		}).Error(err)
	}

	return
}

// ProductsFind func
func (r *ProductsRepository) ProductsFind(ctx context.Context, Limit int, Offset int, Condition map[string]interface{}) (Products []*entities.Products, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("products:%d:%d", Limit, Offset)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache products",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.Select("*").From("products")

		for key, val := range Condition {
			Query.Where(key+" = ?", val)
		}

		_, err = Query.
			Limit(uint64(Limit)).
			Offset(uint64(Offset)).
			LoadContext(ctx, &Products)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query products find",
			}).Error(err)
			return
		}

		// Set Cache
		ProductsJSON, _ := json.Marshal(Products)
		err = r.Redis.Client().Set(CacheKey, ProductsJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for products",
			}).Error(err)
		}
		return Products, nil
	}

	_ = json.Unmarshal([]byte(Value), &Products)
	return Products, nil
}

// ProductsFindOneByID func
func (r *ProductsRepository) ProductsFindOneByID(ctx context.Context, ID int) (Products *entities.Products, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("products:id:%d", ID)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache products by ID",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("products").
			Where("id = ?", ID)

		_, err = Query.LoadContext(ctx, &Products)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query products find by id",
			}).Error(err)
			return
		}

		// Set Cache
		ProductsJSON, _ := json.Marshal(Products)
		err = r.Redis.Client().Set(CacheKey, ProductsJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for products",
			}).Error(err)
		}
		return Products, nil
	}

	_ = json.Unmarshal([]byte(Value), &Products)
	return Products, nil
}

// ProductsStore func
func (r *ProductsRepository) ProductsStore(ctx context.Context, db *dbr.Tx, Products *entities.Products) (ID int, err error) {
	if err = db.InsertInto("products").
		Columns(
			"name",
			"price",
			"qty",
			"status",
			"created_at",
			"updated_at",
		).
		Record(Products).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store products",
		}).Error(err)
	}

	return
}

// ProductsLogStore func
func (r *ProductsRepository) ProductsLogStore(ctx context.Context, db *dbr.Tx, ProductsLog *entities.ProductsLog) (ID int, err error) {
	if err = db.InsertInto("products_log").
		Columns(
			"product_id",
			"user_id",
			"name",
			"price",
			"qty",
			"status",
			"event",
			"created_at",
			"updated_at",
		).
		Record(ProductsLog).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store products log",
		}).Error(err)
	}

	return
}

// ProductsUpdate func
func (r *ProductsRepository) ProductsUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error) {
	_, err = db.Update("products").
		Where("id = ?", ID).
		SetMap(Payload).
		ExecContext(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when update products",
		}).Error(err)
	}

	return
}
