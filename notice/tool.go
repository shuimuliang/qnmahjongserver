package notice

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

// ToolInitAccount ...
func (t *Notice) ToolInitAccount(args *int32, reply *([]string)) error {
	cache.InitAccount()
	return nil
}

// ToolInitCost ...
func (t *Notice) ToolInitCost(args *int32, reply *([]string)) error {
	cache.InitCost()
	return nil
}

// ToolInitGame ...
func (t *Notice) ToolInitGame(args *int32, reply *([]string)) error {
	cache.InitGame()
	return nil
}

// ToolInitModule ...
func (t *Notice) ToolInitModule(args *int32, reply *([]string)) error {
	cache.InitModule()
	return nil
}

// ToolInitPermission ...
func (t *Notice) ToolInitPermission(args *int32, reply *([]string)) error {
	cache.InitPermission()
	return nil
}

// ToolInitRole ...
func (t *Notice) ToolInitRole(args *int32, reply *([]string)) error {
	cache.InitRole()
	return nil
}

// ToolInitShop ...
func (t *Notice) ToolInitShop(args *int32, reply *([]string)) error {
	cache.InitShop()
	return nil
}

// ToolInitAccount ...
func ToolInitAccount() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitAccount", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitCost ...
func ToolInitCost() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitCost", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitGame ...
func ToolInitGame() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitGame", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitModule ...
func ToolInitModule() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitModule", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitPermission ...
func ToolInitPermission() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitPermission", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitRole ...
func ToolInitRole() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitRole", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}

// ToolInitShop ...
func ToolInitShop() {
	client, err := rpc.DialHTTP("tcp", ":5013")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.ToolInitShop", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}
