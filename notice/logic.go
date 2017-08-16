package notice

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

// LogicInitCost ...
func (t *Notice) LogicInitCost(args *int32, reply *([]string)) error {
	cache.InitCost()
	return nil
}

// LogicInitGame ...
func (t *Notice) LogicInitGame(args *int32, reply *([]string)) error {
	cache.InitGame()
	return nil
}

// LogicInitShop ...
func (t *Notice) LogicInitShop(args *int32, reply *([]string)) error {
	cache.InitShop()
	return nil
}

// LogicHandleOrder ...
func (t *Notice) LogicHandleOrder(orderID *string, reply *([]string)) error {
	cache.HandleOrderApply(*orderID)
	return nil
}

// LogicInitCost ...
func LogicInitCost() {
	client, err := rpc.DialHTTP("tcp", ":5012")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.LogicInitCost", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// LogicInitGame ...
func LogicInitGame() {
	client, err := rpc.DialHTTP("tcp", ":5012")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.LogicInitGame", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// LogicInitShop ...
func LogicInitShop() {
	client, err := rpc.DialHTTP("tcp", ":5012")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.LogicInitShop", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// LogicHandleOrder ...
func LogicHandleOrder(orderID string) {
	client, err := rpc.DialHTTP("tcp", ":5012")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := orderID
	reply := make([]string, 10)
	err = client.Call("Notice.LogicHandleOrder", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}
