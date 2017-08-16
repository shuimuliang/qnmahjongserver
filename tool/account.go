package tool

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"
	"qnmahjong/util"

	log "github.com/Sirupsen/logrus"
)

func handleAccountsManage(accounts []*dao.Account) error {
	for _, account := range accounts {
		account.Password = util.Sha1Password(account.Password)
		daoAccount := cache.GetGMTAccount(account.IndexID)
		if daoAccount == nil {
			createAccount(account)
			continue
		}

		account.SetExist(true)
		updateAccount(account)
	}

	notice.ToolInitAccount()
	return nil
}

func createAccount(account *dao.Account) {
	err := account.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertAccount)
	}
}

func updateAccount(account *dao.Account) {
	err := account.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateAccount)
	}
}

func deleteAccount(account *dao.Account) {
	err := account.Delete(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDeleteAccount)
	}
}
