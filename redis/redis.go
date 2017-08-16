package redis

import (
	"fmt"
	"qnmahjong/def"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
)

// Pool is redis pool
var (
	Pool *redis.Pool
)

// Start redis client
func Start() {
	Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", viper.GetString("redis_address"))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error(def.ErrDialRedis)
			}
			return
		},
	}
}

// Shutdown redis client
func Shutdown() {
	if Pool != nil {
		err := Pool.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error(def.ErrCloseRedis)
		}
		fmt.Println("redis shut down")
	}
}
