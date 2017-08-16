package tool

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"

	log "github.com/Sirupsen/logrus"
)

func handleModulesManage(modules []*dao.Module) error {
	for _, module := range modules {
		daoModule := cache.GetGMTModule(module.IndexID)
		if daoModule == nil {
			createModule(module)
			continue
		}

		module.SetExist(true)
		updateModule(module)
	}

	notice.ToolInitModule()
	return nil
}

func createModule(module *dao.Module) {
	err := module.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertModule)
	}
}

func updateModule(module *dao.Module) {
	err := module.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateModule)
	}
}

func deleteModule(module *dao.Module) {
	err := module.Delete(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDeleteModule)
	}
}
