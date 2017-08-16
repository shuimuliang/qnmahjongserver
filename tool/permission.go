package tool

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"

	log "github.com/Sirupsen/logrus"
)

func handlePermissionsManage(permissions []*dao.Permission) error {
	for _, permission := range permissions {
		daoPermission := cache.GetGMTPermission(permission.IndexID)
		if daoPermission == nil {
			createPermission(permission)
			continue
		}

		permission.SetExist(true)
		updatePermission(permission)
	}

	notice.ToolInitPermission()
	return nil
}

func createPermission(permission *dao.Permission) {
	err := permission.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertPermission)
	}
}

func updatePermission(permission *dao.Permission) {
	err := permission.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdatePermission)
	}
}

func deletePermission(permission *dao.Permission) {
	err := permission.Delete(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDeletePermission)
	}
}
