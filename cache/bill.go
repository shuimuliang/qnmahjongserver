package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"sort"
	"sync"
	"time"
)

// AgBillMap hold agBill info(start_time->ag_id->agBill)
type AgBillMap struct {
	sync.RWMutex
	mmap map[string]map[int32]*dao.AgBill
	smap map[int32]*dao.AgBill
}

var (
	agBillCache AgBillMap
)

func init() {
	agBillCache = AgBillMap{}
}

// InitAgBill read sale ag_bill info from db
func InitAgBill() {
	agBillCache.Lock()
	defer agBillCache.Unlock()

	startTimes, err := dao.StartTimesFromAgBill(db.Pool)
	if err != nil {
		return
	}

	agBillCache.mmap = make(map[string]map[int32]*dao.AgBill, len(startTimes))
	agBillCache.smap = make(map[int32]*dao.AgBill, len(startTimes))
	for _, startTime := range startTimes {
		agBills, err := dao.AgBillsByStartTime(db.Pool, startTime)
		if err != nil {
			continue
		}

		start := startTime.Format("2006-01-02")
		agBillCache.mmap[start] = make(map[int32]*dao.AgBill, len(agBills))
		for _, agBill := range agBills {
			agBillCache.mmap[start][agBill.AgID] = agBill
			agBillCache.smap[agBill.IndexID] = agBill
		}
	}
}

func GetQueryBalanceRD(agID int32) []*dao.AgBill {
	agBillCache.RLock()
	defer agBillCache.RUnlock()

	rds := make([]*dao.AgBill, 0)
	for _, agBillMap := range agBillCache.mmap {
		agBill, ok := agBillMap[agID]
		if ok {
			rds = append(rds, agBill)
		}
	}

	sort.Slice(rds, func(i, j int) bool {
		return rds[i].StartTime.Unix() >= rds[j].StartTime.Unix()
	})

	return rds
}

func GetVerifyBalanceRD() []*VerifyBalanceRD {
	agBillCache.RLock()
	defer agBillCache.RUnlock()

	rds := make([]*VerifyBalanceRD, 0)
	for start := range agBillCache.mmap {
		startTime, _ := time.ParseInLocation("2006-01-02", start, time.Local)
		endTime := startTime.AddDate(0, 0, 7)
		end := endTime.Format("2006-01-02")
		rd := &VerifyBalanceRD{
			Start: start,
			End:   end,
		}
		rds = append(rds, rd)
	}

	sort.Slice(rds, func(i, j int) bool {
		return rds[i].Start >= rds[j].Start
	})

	return rds
}

func GetVerifyBalanceDetailRD(start string) []*VerifyBalanceDetailRD {
	agBillCache.RLock()
	defer agBillCache.RUnlock()
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	agBillMap, ok := agBillCache.mmap[start]
	if !ok {
		return nil
	}

	rds := make([]*VerifyBalanceDetailRD, 0)
	for agID, agBill := range agBillMap {
		agAgent, ok := agAgentCache.mmap[agID]
		if !ok {
			continue
		}

		rds = append(rds, &VerifyBalanceDetailRD{
			AgAgent: agAgent,
			AgBill:  agBill,
		})
	}

	sort.Slice(rds, func(i, j int) bool {
		return rds[i].AgID < rds[j].AgID
	})

	return rds
}

func UpdateAgBillStatus(agID, indexID, status int32) bool {
	agBillCache.RLock()
	defer agBillCache.RUnlock()
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return false
	}

	agBill, ok := agBillCache.smap[indexID]
	if !ok {
		return false
	}

	agAgent.TotalBalance += agBill.LastWeekDakuan
	agBill.Delflag = status

	agAgent.AgAccount.Update(db.Pool)
	agBill.Update(db.Pool)
	return true
}
