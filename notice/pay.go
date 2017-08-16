package notice

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
)

// PayInitShop ...
func (t *Notice) PayInitShop(args *int32, reply *([]string)) error {
	cache.InitShop()
	return nil
}

// PayInitShop ...
func PayInitShop() {
	client, err := rpc.DialHTTP("tcp", ":5014")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDialNoticeRPC)
		return
	}

	args := 0
	reply := make([]string, 10)
	err = client.Call("Notice.PayInitShop", &args, &reply)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrCallNoticeRPC)
	}
	client.Close()
}
