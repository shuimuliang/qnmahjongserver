package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"sync"
)

// CostMap hold mj createroom cost info(mjtype->round->cost)
type CostMap struct {
	sync.RWMutex
	mmap map[int32]map[int32]*dao.Cost
	gmap map[int32]*dao.Cost
}

var (
	costCache CostMap
)

func init() {
	costCache = CostMap{}
}

// InitCost read mj createroom info from db
func InitCost() {
	costCache.Lock()
	defer costCache.Unlock()

	mjTypes, err := dao.MjTypesFromCost(db.Pool)
	if err != nil {
		return
	}

	costCache.mmap = make(map[int32]map[int32]*dao.Cost, len(mjTypes))
	costCache.gmap = make(map[int32]*dao.Cost, len(mjTypes))
	for _, mjType := range mjTypes {
		costs, err := dao.CostsByMjType(db.Pool, mjType)
		if err != nil {
			continue
		}

		costCache.mmap[mjType] = make(map[int32]*dao.Cost, len(costs))
		for _, cost := range costs {
			costCache.mmap[mjType][cost.Rounds] = cost
			costCache.gmap[cost.IndexID] = cost
		}
	}
}

// GetJushus get jushus by mj types
func GetJushus(mjTypes []int32) []*pf.Jushu {
	costCache.RLock()
	defer costCache.RUnlock()

	var jushus []*pf.Jushu
	for _, mjType := range mjTypes {
		costs, ok := costCache.mmap[mjType]
		if !ok {
			continue
		}

		for _, cost := range costs {
			jushus = append(jushus, &pf.Jushu{
				Jushu:  cost.Rounds,
				Coins:  cost.Coins,
				Cards:  cost.Cards,
				MjType: cost.MjType,
				MjDesc: cost.MjDesc,
			})
		}
	}
	return jushus
}

// GetCosts get costs by mj types and rounds
func GetCosts(mjType, rounds, costType int32) (int32, bool) {
	costCache.RLock()
	defer costCache.RUnlock()

	costs, ok := costCache.mmap[mjType]
	if !ok {
		return 0, false
	}

	cost, ok := costs[rounds]
	if !ok {
		return 0, false
	}

	if costType == def.CostCard {
		return cost.Cards, true
	}

	return cost.Coins, true
}

// GetGMTCost get cost by index_id
func GetGMTCost(IndexID int32) *dao.Cost {
	costCache.RLock()
	defer costCache.RUnlock()

	return costCache.gmap[IndexID]
}
