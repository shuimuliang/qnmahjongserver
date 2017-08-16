package util

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
)

// Stack print panic msg
func Stack() {
	if err := recover(); err != nil {
		log.WithFields(log.Fields{
			"stack": stack(),
		}).Error(err)
	}
}

func stack() string {
	buf := make([]byte, 1024*10)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
