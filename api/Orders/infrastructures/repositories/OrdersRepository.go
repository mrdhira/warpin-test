package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
	dbr "github.com/gocraft/dbr/v2"
	"github.com/mrdhira/warpin-test/api/Orders/entities"
	"github.com/mrdhira/warpin-test/api/Orders/infrastructures/database"
	log "github.com/sirupsen/logrus"
)

// IOrdersRepository interface
type IOrdersRepository interface {
	Tx() (tx *dbr.Tx, err error)
	OrdersFind(ctx context.Context, Limit int, Offset int, Condition map[string]interface{}) (Orders []*entities.Orders, err error)
	OrdersFindByUserID(ctx context.Context, Limit int, Offset int, UserID int) (Orders []*entities.Orders, err error)
	OrdersFindByID(ctx context.Context, ID int) (Orders *entities.Orders, err error)
	OrdersStore(ctx context.Context, db *dbr.Tx, Orders *entities.Orders) (ID int, err error)
	OrdersLogStore(ctx context.Context, db *dbr.Tx, OrdersLog *entities.OrdersLog) (ID int, err error)
	OrdersUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error)
}

// OrdersRepository struct
type OrdersRepository struct {
	PG    database.IPostgresConnection
	Redis database.IRedisConnection
}

// Tx func to create new transaction
func (r *OrdersRepository) Tx() (tx *dbr.Tx, err error) {
	db := r.PG.PostgresTrade()

	tx, err = db.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when begin transaction in postgres",
		}).Error(err)
	}

	return
}

// OrdersFind func
func (r *OrdersRepository) OrdersFind(ctx context.Context, Limit int, Offset int, Condition map[string]interface{}) (Orders []*entities.Orders, err error) {
	db := r.PG.PostgresTrade()

	Query := db.Select("*").From("orders")

	for key, val := range Condition {
		Query.Where(key+" = ?", val)
	}

	_, err = Query.
		Limit(uint64(Limit)).
		Offset(uint64(Offset)).
		LoadContext(ctx, &Orders)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when query orders find",
		}).Error(err)
		return
	}

	return
}

// OrdersFindByUserID func
func (r *OrdersRepository) OrdersFindByUserID(ctx context.Context, Limit int, Offset int, UserID int) (Orders []*entities.Orders, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("orders:users:%d:%d:%d", UserID, Limit, Offset)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache orders by user id",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("orders").
			Where("user_id = ?", UserID)

		_, err = Query.
			Limit(uint64(Limit)).
			Offset(uint64(Offset)).
			LoadContext(ctx, &Orders)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query orders find by user id",
			}).Error(err)
			return
		}

		// Set Cache
		OrdersJSON, _ := json.Marshal(Orders)
		err = r.Redis.Client().Set(CacheKey, OrdersJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for orders by user id",
			}).Error(err)
		}
		return Orders, nil
	}

	_ = json.Unmarshal([]byte(Value), &Orders)
	return Orders, nil
}

// OrdersFindByID func
func (r *OrdersRepository) OrdersFindByID(ctx context.Context, ID int) (Orders *entities.Orders, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("orders:id:%d", ID)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache orders by id",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("orders").
			Where("id = ?", ID)

		_, err = Query.
			Limit(1).
			LoadContext(ctx, &Orders)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query orders find by id",
			}).Error(err)
			return
		}

		// Set Cache
		OrdersJSON, _ := json.Marshal(Orders)
		err = r.Redis.Client().Set(CacheKey, OrdersJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for orders by id",
			}).Error(err)
		}
		return Orders, nil
	}

	_ = json.Unmarshal([]byte(Value), &Orders)
	return Orders, nil
}

// OrdersStore func
func (r *OrdersRepository) OrdersStore(ctx context.Context, db *dbr.Tx, Orders *entities.Orders) (ID int, err error) {
	if err = db.InsertInto("orders").
		Columns(
			"user_id",
			"product_id",
			"product_name",
			"price",
			"qty",
			"total_price",
			"status",
			"created_at",
			"updated_at",
		).
		Record(Orders).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store orders",
		}).Error(err)
	}

	return
}

// OrdersLogStore func
func (r *OrdersRepository) OrdersLogStore(ctx context.Context, db *dbr.Tx, OrdersLog *entities.OrdersLog) (ID int, err error) {
	if err = db.InsertInto("orders_log").
		Columns(
			"user_id",
			"product_id",
			"product_name",
			"price",
			"qty",
			"total_price",
			"status",
			"event",
			"admin_id",
			"created_at",
			"updated_at",
		).
		Record(OrdersLog).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store orders log",
		}).Error(err)
	}

	return
}

// OrdersUpdate func
func (r *OrdersRepository) OrdersUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error) {
	_, err = db.Update("orders").
		Where("id = ?", ID).
		SetMap(Payload).
		ExecContext(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when update orders",
		}).Error(err)
	}

	return
}
