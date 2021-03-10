package database

import (
	"strconv"

	dbr "github.com/gocraft/dbr/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "go.elastic.co/apm/module/apmsql/pq"
)

// IPostgresConnection interface
type IPostgresConnection interface {
	PostgresTrade() *dbr.Session
}

// PostgresConnection struct
type PostgresConnection struct{}

// Initialize Variable
var (
	PGConnection *dbr.Connection
)

// PostgresTrade func
func (p *PostgresConnection) PostgresTrade() *dbr.Session {
	if PGConnection == nil {
		Driver := viper.GetString("ordersServices.database.driver")
		DSN := viper.GetString("ordersServices.database.dsn")
		MaxIdle, _ := strconv.Atoi(viper.GetString("ordersServices.database.max_idle"))
		MaxConn, _ := strconv.Atoi(viper.GetString("ordersServices.database.max_conn"))

		var err error

		PGConnection, err = dbr.Open(Driver, DSN, nil)
		if err != nil {
			log.WithFields(log.Fields{
				"event": "error when create sql connection",
			}).Error(err)
		}
		PGConnection.SetMaxIdleConns(MaxIdle)
		PGConnection.SetMaxOpenConns(MaxConn)
	}

	Session := PGConnection.NewSession(nil)
	return Session
}
