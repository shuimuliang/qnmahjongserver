package db

import (
	"database/sql"
	"fmt"
	"qnmahjong/def"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// Pool is database pool
var (
	Pool *sql.DB
)

// Start db client
func Start() {
	var err error
	Pool, err = sql.Open("mysql", viper.GetString("mysql_address"))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrOpenDb)
	}
	if err = Pool.Ping(); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrOpenDb)
		panic(err)
	}
}

// Shutdown db client
func Shutdown() {
	if Pool != nil {
		err := Pool.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrCloseDb)
		}
		fmt.Println("db shut down")
	}
}
