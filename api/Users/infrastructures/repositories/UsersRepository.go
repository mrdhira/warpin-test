package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
	dbr "github.com/gocraft/dbr/v2"
	"github.com/mrdhira/warpin-test/api/Users/entities"
	"github.com/mrdhira/warpin-test/api/Users/infrastructures/database"
	log "github.com/sirupsen/logrus"
)

// IUsersRepository interface
type IUsersRepository interface {
	Tx() (tx *dbr.Tx, err error)
	UsersFindOne(ctx context.Context, Condition map[string]interface{}) (Users *entities.Users, err error)
	UsersFindByID(ctx context.Context, ID int) (Users *entities.Users, err error)
	UsersFindByEmail(ctx context.Context, Email string) (Users *entities.Users, err error)
	UsersStore(ctx context.Context, db *dbr.Tx, Users *entities.Users) (ID int, err error)
	UsersLogStore(ctx context.Context, db *dbr.Tx, UsersLog *entities.UsersLog) (ID int, err error)
	// UsersEventLogStore(ctx context.Context, db *dbr.Tx, UsersEventLog *entities.UsersEventLog) (ID int, err error)
	ProfileByID(ctx context.Context, ID int) (Profile *entities.Profile, err error)
	UsersUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error)
}

// UsersRepository struct
type UsersRepository struct {
	PG    database.IPostgresConnection
	Redis database.IRedisConnection
}

// Tx func to create new transaction
func (r *UsersRepository) Tx() (tx *dbr.Tx, err error) {
	db := r.PG.PostgresTrade()

	tx, err = db.Begin()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when begin transaction in postgres",
		}).Error(err)
	}

	return
}

// UsersFindOne func
func (r *UsersRepository) UsersFindOne(ctx context.Context, Condition map[string]interface{}) (Users *entities.Users, err error) {
	db := r.PG.PostgresTrade()

	Query := db.Select("*").From("users")

	for key, val := range Condition {
		Query.Where(key+" = ?", val)
	}

	_, err = Query.Limit(1).LoadContext(ctx, &Users)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when query users find one",
		}).Error(err)
		return
	}

	return
}

// UsersFindByID func
func (r *UsersRepository) UsersFindByID(ctx context.Context, ID int) (Users *entities.Users, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("users:id:%d", ID)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache users by ID",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("users").
			Where("id = ?", ID)

		_, err = Query.LoadContext(ctx, &Users)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query users find by id",
			}).Error(err)
			return
		}

		// Set Cache
		UsersJSON, _ := json.Marshal(Users)
		err = r.Redis.Client().Set(CacheKey, UsersJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for users",
			}).Error(err)
		}
		return Users, nil
	}

	_ = json.Unmarshal([]byte(Value), &Users)
	return Users, nil
}

// UsersFindByEmail func
func (r *UsersRepository) UsersFindByEmail(ctx context.Context, Email string) (Users *entities.Users, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("users:email:%s", Email)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache users by email",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("users").
			Where("email = ?", Email)

		_, err = Query.LoadContext(ctx, &Users)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query users find by email",
			}).Error(err)
			return
		}

		// Set Cache
		UsersJSON, _ := json.Marshal(Users)
		err = r.Redis.Client().Set(CacheKey, UsersJSON, time.Second*10).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for users",
			}).Error(err)
		}
		return
	}

	_ = json.Unmarshal([]byte(Value), &Users)
	return Users, nil
}

// UsersStore func
func (r *UsersRepository) UsersStore(ctx context.Context, db *dbr.Tx, Users *entities.Users) (ID int, err error) {
	if err = db.InsertInto("users").
		Columns(
			"email",
			"phone_number",
			"full_name",
			"gender",
			"role",
			"password",
			"created_at",
			"updated_at",
		).
		Record(Users).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store users",
		}).Error(err)
	}

	return
}

// UsersLogStore func
func (r *UsersRepository) UsersLogStore(ctx context.Context, db *dbr.Tx, UsersLog *entities.UsersLog) (ID int, err error) {
	if err = db.InsertInto("users_log").
		Columns(
			"user_id",
			"email",
			"phone_number",
			"full_name",
			"gender",
			"role",
			"password",
			"created_at",
			"updated_at",
		).
		Record(UsersLog).
		Returning("id").
		LoadContext(ctx, &ID); err != nil {
		log.WithFields(log.Fields{
			"event": "error when store users log",
		}).Error(err)
	}

	return
}

// ProfileByID func
func (r *UsersRepository) ProfileByID(ctx context.Context, ID int) (Profile *entities.Profile, err error) {
	db := r.PG.PostgresTrade()

	// Check Cache
	CacheKey := fmt.Sprintf("users:profile:id:%d", ID)
	Value, err := r.Redis.Client().Get(CacheKey).Result()
	if err != redis.Nil && err != nil {
		log.WithFields(log.Fields{
			"event": "error when get cache profile by ID",
		}).Error(err)
	}

	if err == redis.Nil {
		Query := db.
			Select("*").
			From("users").
			Where("id = ?", ID)

		_, err = Query.LoadContext(ctx, &Profile)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when query profile find by id",
			}).Error(err)
			return
		}

		// Set Cache
		ProfileJSON, _ := json.Marshal(Profile)
		err = r.Redis.Client().Set(CacheKey, ProfileJSON, time.Minute*15).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when set cache for Profile",
			}).Error(err)
		}
	}

	_ = json.Unmarshal([]byte(Value), &Profile)
	return
}

// UsersUpdate func
func (r *UsersRepository) UsersUpdate(ctx context.Context, db *dbr.Tx, ID int, Payload map[string]interface{}) (err error) {
	_, err = db.Update("users").
		Where("id = ?", ID).
		SetMap(Payload).
		ExecContext(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when update users",
		}).Error(err)
	}

	return
}
