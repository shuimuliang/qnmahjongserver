package tool

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"

	log "github.com/Sirupsen/logrus"
)

func handleRolesManage(roles []*dao.Role) error {
	for _, role := range roles {
		daoRole := cache.GetGMTRole(role.IndexID)
		if daoRole == nil {
			createRole(role)
			continue
		}

		role.SetExist(true)
		updateRole(role)
	}

	notice.ToolInitRole()
	return nil
}

func createRole(role *dao.Role) {
	err := role.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertRole)
	}
}

func updateRole(role *dao.Role) {
	err := role.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateRole)
	}
}

func deleteRole(role *dao.Role) {
	err := role.Delete(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDeleteRole)
	}
}
