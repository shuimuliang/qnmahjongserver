package notice

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

// SaleInitAgAuth ...
func (t *Notice) SaleInitAgAuth(args *int32, reply *([]string)) error {
	cache.InitAgAuth()
	return nil
}

// SaleInitAgAccount ...
func (t *Notice) SaleInitAgAccount(args *int32, reply *([]string)) error {
	cache.InitAgAccount()
	return nil
}

// SaleInitAgBill ...
func (t *Notice) SaleInitAgBill(args *int32, reply *([]string)) error {
	cache.InitAgBill()
	return nil
}

// SaleInitAgAuth ...
func SaleInitAgAuth() {
	client, err := rpc.DialHTTP("tcp", ":5015")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.SaleInitAgAuth", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// SaleInitAgAccount ...
func SaleInitAgAccount() {
	client, err := rpc.DialHTTP("tcp", ":5015")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.SaleInitAgAccount", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// SaleInitAgBill ...
func SaleInitAgBill() {
	client, err := rpc.DialHTTP("tcp", ":5015")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.SaleInitAgBill", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}
