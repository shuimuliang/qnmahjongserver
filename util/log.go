package util

import (
	log "github.com/Sirupsen/logrus"
)

// LogSend ...
func LogSend(msgID, id, roomID int32, send interface{}, info string) {
	log.WithFields(log.Fields{
		"msgid":    msgID,
		"playerid": id,
		"roomid":   roomID,
		"send":     send,
	}).Info(info)
}

// LogRecv ...
func LogRecv(msgID, id, roomID int32, recv interface{}, info string) {
	log.WithFields(log.Fields{
		"msgid":    msgID,
		"playerid": id,
		"roomid":   roomID,
		"recv":     recv,
	}).Info(info)
}

// LogError ...
func LogError(err error, key string, value interface{}, id int32, customErr error) {
	log.WithFields(log.Fields{
		"error":    err,
		key:        value,
		"playerid": id,
	}).Error(customErr)
}
