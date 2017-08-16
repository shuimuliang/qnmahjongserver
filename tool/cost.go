package tool

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/notice"

	log "github.com/Sirupsen/logrus"
)

func handleCostsManage(costs []*dao.Cost) error {
	for _, cost := range costs {
		daoCost := cache.GetGMTCost(cost.IndexID)
		if daoCost == nil {
			createCost(cost)
			continue
		}

		cost.SetExist(true)
		updateCost(cost)
	}

	notice.ToolInitCost()
	return nil
}

func createCost(cost *dao.Cost) {
	err := cost.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertCost)
	}
}

func updateCost(cost *dao.Cost) {
	err := cost.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateCost)
	}
}

func deleteCost(cost *dao.Cost) {
	err := cost.Delete(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrDeleteCost)
	}
}
