package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

// AgAuthMap hold sale ag_auth info(ag_upper_id->ag_id->ag_auth)
type AgAuthMap struct {
	sync.RWMutex
	mmap map[int32]map[int32]*dao.AgAuth
}

var (
	agAuthCache AgAuthMap
)

func init() {
	agAuthCache = AgAuthMap{}
}

// InitAgAuth read sale ag_auth info from db
func InitAgAuth() {
	agAuthCache.Lock()
	defer agAuthCache.Unlock()

	agIDs, err := dao.AgUpperIDsFromAgAuth(db.Pool)
	if err != nil {
		return
	}

	agAuthCache.mmap = make(map[int32]map[int32]*dao.AgAuth, len(agIDs))
	for _, agID := range agIDs {
		agAuths, err := dao.AgAuthsByAgUpperID(db.Pool, agID)
		if err != nil {
			continue
		}

		agAuthCache.mmap[agID] = make(map[int32]*dao.AgAuth, len(agAuths))
		for _, agAuth := range agAuths {
			agAuthCache.mmap[agID][agAuth.AgID] = agAuth
		}
	}
}

// CreateAgAuth create sale agAuth
func CreateAgAuth(agUpperID, agID, agLevel int32) bool {
	agAuthCache.RLock()
	defer agAuthCache.RUnlock()

	agAuths, ok := agAuthCache.mmap[agUpperID]
	if ok {
		_, ok := agAuths[agID]
		if ok {
			return false
		}
	}

	agAuth := &dao.AgAuth{
		AgUpperID:  agUpperID,
		AgID:       agID,
		AgLevel:    agLevel,
		CreateTime: time.Now(),
	}
	err := agAuth.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertAgAuth)
		return false
	}

	return true
}

// CheckAgAuth create sale agAuth
func CheckAgAuth(agUpperID, agID int32) *dao.AgAuth {
	agAuthCache.RLock()
	defer agAuthCache.RUnlock()

	agAuths, ok := agAuthCache.mmap[agUpperID]
	if !ok {
		return nil
	}

	agAuth, ok := agAuths[agID]
	if !ok {
		return nil
	}

	// 这里返回是有问题的
	return agAuth
}
