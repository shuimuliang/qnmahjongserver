package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"sync"
)

// AccountMap hold tool account info(email->account)
type AccountMap struct {
	sync.RWMutex
	mmap map[string]*dao.Account
	gmap map[int32]*dao.Account
}

var (
	accountCache AccountMap
)

func init() {
	accountCache = AccountMap{}
}

// InitAccount read tool account info from db
func InitAccount() {
	accountCache.Lock()
	defer accountCache.Unlock()

	emails, err := dao.EmailsFromAccount(db.Pool)
	if err != nil {
		return
	}

	accountCache.mmap = make(map[string]*dao.Account, len(emails))
	accountCache.gmap = make(map[int32]*dao.Account, len(emails))
	for _, email := range emails {
		daoAccount, err := dao.AccountByEmail(db.Pool, email)
		if err != nil {
			continue
		}

		accountCache.mmap[email] = daoAccount
		accountCache.gmap[daoAccount.IndexID] = daoAccount
	}
}

// CheckAccount check tool account password
func CheckAccount(email, password string) bool {
	accountCache.RLock()
	defer accountCache.RUnlock()

	account, ok := accountCache.mmap[email]
	if !ok {
		return false
	}

	return account.Password == password
}

// GetGMTAccount get account by index_id
func GetGMTAccount(IndexID int32) *dao.Account {
	accountCache.RLock()
	defer accountCache.RUnlock()

	return accountCache.gmap[IndexID]
}
