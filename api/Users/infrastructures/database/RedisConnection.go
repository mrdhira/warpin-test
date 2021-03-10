package database

import (
	"strconv"

	redis "github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IRedisConnection interface
type IRedisConnection interface {
	Client() *redis.Client
}

// RedisConnection struct
type RedisConnection struct{}

// Client Func
func (r *RedisConnection) Client() *redis.Client {
	Address := viper.GetString("usersServices.redis.address")
	Password := viper.GetString("usersServices.redis.password")
	DB, _ := strconv.Atoi(viper.GetString("usersServices.redis.db"))

	Client := redis.NewClient(&redis.Options{
		Addr:     Address,
		Password: Password,
		DB:       DB,
	})

	_, err := Client.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when try to ping redis",
		}).Error(err)
	}

	return Client
}
