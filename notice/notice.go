package notice

import (
	"qnmahjong/def"
	"net"
	"net/http"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

type Notice int32

func StartLoginNotice() {
	notice := new(Notice)

	rpc.Register(notice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5011")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Error(def.ErrStartLoginNotice)
	}

	go http.Serve(l, nil)
}

func StartLogicNotice() {
	notice := new(Notice)

	rpc.Register(notice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5012")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Error(def.ErrStartLogicNotice)
	}

	go http.Serve(l, nil)
}

func StartToolNotice() {
	notice := new(Notice)

	rpc.Register(notice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5013")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Error(def.ErrStartToolNotice)
	}

	go http.Serve(l, nil)
}

func StartPayNotice() {
	notice := new(Notice)

	rpc.Register(notice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5014")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Error(def.ErrStartPayNotice)
	}

	go http.Serve(l, nil)
}

func StartSaleNotice() {
	notice := new(Notice)

	rpc.Register(notice)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":5015")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Error(def.ErrStartPayNotice)
	}

	go http.Serve(l, nil)
}
