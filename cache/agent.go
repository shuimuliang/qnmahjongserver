package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/util"
	"sort"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

type AgAgent struct {
	*dao.AgAccount
	*dao.Player
}

// AgAgentMap hold sale ag_auth info(ag_id->ag_account)
type AgAgentMap struct {
	sync.RWMutex
	mmap map[int32]*AgAgent
}

var (
	agAgentCache AgAgentMap
)

func init() {
	agAgentCache = AgAgentMap{}
}

// InitAgAccount read sale ag_auth info from db
func InitAgAccount() {
	agAgentCache.Lock()
	defer agAgentCache.Unlock()

	agIDs, err := dao.AgIDsFromAgAccount(db.Pool)
	if err != nil {
		return
	}

	agAgentCache.mmap = make(map[int32]*AgAgent, len(agIDs))
	for _, agID := range agIDs {
		agAgent, err := dao.AgAccountByAgID(db.Pool, agID)
		if err != nil {
			continue
		}

		player, err := dao.PlayerByPlayerID(db.Pool, agID)
		if err != nil {
			continue
		}

		agAgentCache.mmap[agID] = &AgAgent{
			AgAccount: agAgent,
			Player:    player,
		}
	}
}

// CheckAgAgent check sale agAgent password
func CheckAgAgent(agID int32, password string) bool {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return false
	}

	return agAgent.Password == password
}

// CreateAgAgent create sale agAgent
func CreateAgAgent(agUpperID, agID int32, password, telephone string) bool {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	_, ok := agAgentCache.mmap[agID]
	if ok {
		return false
	}

	agAuth := CheckAgAuth(agUpperID, agID)
	if agAuth == nil {
		return false
	}

	AgAccount := &dao.AgAccount{
		AgUpperID:  agUpperID,
		AgID:       agID,
		AgLevel:    agAuth.AgLevel,
		Password:   util.Sha1Password(password),
		Telephone:  telephone,
		CreateTime: time.Now(),
	}
	err := AgAccount.Insert(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrInsertAgAccount)
		return false
	}

	return true
}

func getAgLevelRate(agLevel int32) float64 {
	switch agLevel {
	case def.AgLevelOne:
		return def.AgLevelOneRate
	case def.AgLevelTwo:
		return def.AgLevelTwoRate
	default:
		return 0
	}
}

func getFirstBuyRate(agLevel int32) float64 {
	switch agLevel {
	case def.AgLevelOne, def.AgLevelTwo:
		return def.AgFirstBuyRate
	default:
		return 0
	}
}

func getLowAgentsRate(agLevel int32) float64 {
	switch agLevel {
	case def.AgLevelOne, def.AgLevelTwo:
		return def.AgLowAgentsRate
	default:
		return 0
	}
}

func GetSaleIndexRD(agID int32) *SaleIndexRD {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return nil
	}

	agAgents, err := dao.AgAccountsByAgUpperID(db.Pool, agID)
	if err != nil {
		return nil
	}

	players, err := dao.PlayersByHighID(db.Pool, agID)
	if err != nil {
		return nil
	}

	cardCnt, _ := GetTodayRD(agID)
	return &SaleIndexRD{
		AgAgent:    agAgent,
		CardCnt:    cardCnt,
		BalanceCnt: 0,
		AgentCnt:   int32(len(agAgents)),
		PlayerCnt:  int32(len(players)),
	}
}

func GetTodayRD(agID int32) (int32, int32) {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	today := time.Now().Format("2006-01-02")
	agPays, err := dao.AgPaysByAgID(db.Pool, agID)
	if err != nil {
		return 0, 0
	}

	for _, agPay := range agPays {
		if agPay.CreateTime.Format("2006-01-02") == today {
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}
	return cardsCnt, moneyCnt
}

func GetWeekRD(agID int32) (int32, int32) {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	agPays, err := dao.AgPaysByAgID(db.Pool, agID)
	if err != nil {
		return 0, 0
	}

	for _, agPay := range agPays {
		if agPay.Delflag == 0 {
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}
	return cardsCnt, moneyCnt
}

func GetActivityRD(agID int32) (int32, int32) {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	agPays, err := dao.AgPaysByAgID(db.Pool, agID)
	if err != nil {
		return 0, 0
	}

	for _, agPay := range agPays {
		if agPay.Delflag == 0 && agPay.FirstBuyAward == 1 {
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}
	return cardsCnt, moneyCnt
}

func GetLowAgentsRD(agID int32) (int32, int32) {
	cardsCnt := int32(0)
	moneyCnt := int32(0)

	return cardsCnt, moneyCnt
}

func GetBalanceRD(agID int32) (int32, int32) {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	today := time.Now().Format("2006-01-02")
	agPays, err := dao.AgPaysByAgID(db.Pool, agID)
	if err != nil {
		return 0, 0
	}

	for _, agPay := range agPays {
		if agPay.CreateTime.Format("2006-01-02") == today {
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}
	return cardsCnt, moneyCnt
}

func GetUserCardsRD(agID int32, startDate, endDate string) *UserCardsRD {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	orderList := make([]*dao.AgPay, 0)
	start, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, _ := time.ParseInLocation("2006-01-02", endDate, time.Local)
	end = end.AddDate(0, 0, 1)

	agPays, err := dao.AgPaysByAgID(db.Pool, agID)
	if err != nil {
		return nil
	}

	for _, agPay := range agPays {
		if agPay.CreateTime.Unix() >= start.Unix() && agPay.CreateTime.Unix() < end.Unix() {
			orderList = append(orderList, agPay)
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}

	return &UserCardsRD{
		StartDate:   startDate,
		EndDate:     endDate,
		CardsCnt:    cardsCnt,
		MoneyCnt:    moneyCnt,
		OrderLength: int32(len(orderList)),
		OrderList:   orderList,
	}
}

func GetUserBalanceRD(agID int32) *UserBalanceRD {
	agAgentCache.Lock()
	defer agAgentCache.Unlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return nil
	}

	_, todayCnt := GetTodayRD(agID)
	todayBalance := int32(float64(todayCnt) * getAgLevelRate(agAgent.AgLevel))

	_, weekCnt := GetWeekRD(agID)
	weekBalance := int32(float64(weekCnt) * getAgLevelRate(agAgent.AgLevel))

	_, activityCnt := GetWeekRD(agID)
	activityBalance := int32(float64(activityCnt) * 0.1)
	if agAgent.AgLevel == 0 {
		activityBalance = 0
	}

	hongbao := int32(0)
	if agAgent.Hongbao == 0 {
		hongbao = 50000
	}

	return &UserBalanceRD{
		Balance1: todayBalance,
		Balance2: 100,
		Balance3: 0,
		Balance4: weekBalance,
		Balance5: activityBalance,
		Balance6: hongbao,
		Balance7: agAgent.TotalBalance,
	}
}

func GetAgentAuthRD(agID int32) *AgAgent {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	return agAgentCache.mmap[agID]
}

func GetQueryOrderRD(agID, playerID int32, startDate, endDate string) *QueryOrderRD {
	cardsCnt := int32(0)
	moneyCnt := int32(0)
	orderList := make([]*dao.AgPay, 0)
	start, _ := time.ParseInLocation("2006-01-02", startDate, time.Local)
	end, _ := time.ParseInLocation("2006-01-02", endDate, time.Local)

	agPays, err := dao.AgPaysByAgIDCustomerID(db.Pool, agID, playerID)
	if err != nil {
		return nil
	}

	for _, agPay := range agPays {
		if agPay.CreateTime.After(start) && agPay.CreateTime.Before(end) {
			orderList = append(orderList, agPay)
			cardsCnt += agPay.DiamondCnt
			moneyCnt += agPay.MoneyCnt
		}
	}

	return &QueryOrderRD{
		StartDate:   startDate,
		EndDate:     endDate,
		PlayerID:    playerID,
		CardsCnt:    cardsCnt,
		MoneyCnt:    moneyCnt,
		OrderLength: int32(len(orderList)),
		OrderList:   orderList,
	}
}

func GetMyAgentsRD(agID int32) []*AgAgent {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	rds := make([]*AgAgent, 0)
	for _, agAgent := range agAgentCache.mmap {
		if agAgent.AgUpperID == agID {
			rds = append(rds, agAgent)
		}
	}

	sort.Slice(rds, func(i, j int) bool {
		return rds[i].PlayerID < rds[j].PlayerID
	})

	return rds
}

func GetMyPlayersRD(agID int32) []*dao.Player {
	rds, err := dao.PlayersByHighID(db.Pool, agID)
	if err != nil {
		return nil
	}

	sort.Slice(rds, func(i, j int) bool {
		return rds[i].PlayerID < rds[j].PlayerID
	})

	return rds
}

func GetEditProfileRD(agID int32) *AgAgent {
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	return agAgentCache.mmap[agID]
}

func EditAgAgentProfile(agID int32, telephone, realname, weixin, alipay, email string) bool {
	agAgentCache.Lock()
	defer agAgentCache.Unlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return false
	}

	agAccount := agAgent.AgAccount
	agAccount.Telephone = telephone
	agAccount.Realname = realname
	agAccount.Weixin = weixin
	agAccount.Alipay = alipay
	agAccount.Email = email
	err := agAccount.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateAgAccount)
		return false
	}

	return true
}

func EditAgAgentPWD(agID int32, password string) bool {
	agAgentCache.Lock()
	defer agAgentCache.Unlock()

	agAgent, ok := agAgentCache.mmap[agID]
	if !ok {
		return false
	}

	agAccount := agAgent.AgAccount
	oldPassword := agAccount.Password
	agAccount.Password = password
	err := agAccount.Update(db.Pool)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error(def.ErrUpdateAgAccount)
		agAccount.Password = oldPassword
		return false
	}

	return true
}

func CashSettlement() {
	agBillCache.RLock()
	defer agBillCache.RUnlock()
	agAgentCache.RLock()
	defer agAgentCache.RUnlock()

	now := time.Now()
	end, _ := time.ParseInLocation("2006-01-02", now.Format("2006-01-02"), time.Local)
	start := end.AddDate(0, 0, -7)
	last := start.AddDate(0, 0, -7).Format("2006-01-02")

	for agID, agAgent := range agAgentCache.mmap {
		lastWeekLeft := int32(0)
		agBillMap, ok := agBillCache.mmap[last]
		if ok {
			agBill, ok := agBillMap[agID]
			if ok {
				if agBill.Delflag != 1 {
					lastWeekLeft = agBill.LastWeekDakuan
				}
			}
		}

		_, lowAgentsCnt := GetLowAgentsRD(agID)
		_, weekCnt := GetWeekRD(agID)
		_, activityCnt := GetActivityRD(agID)
		lowAgentsAward := int32(float64(lowAgentsCnt) * getLowAgentsRate(agAgent.AgLevel))
		cardsAward := int32(float64(weekCnt) * getAgLevelRate(agAgent.AgLevel))
		firstBuyAward := int32(float64(activityCnt) * getFirstBuyRate(agAgent.AgLevel))
		honbao := int32(0)
		lastWeekBalance := lowAgentsAward + cardsAward + firstBuyAward + honbao
		lastWeekDakuan := lastWeekLeft + lastWeekBalance

		agBill := dao.AgBill{
			AgID:            agID,
			LastWeekLeft:    lastWeekLeft,
			LastWeekBalance: lastWeekBalance,
			LastWeekDakuan:  lastWeekDakuan,
			LowAgentsAward:  lowAgentsAward,
			CardsAward:      cardsAward,
			FirstBuyAward:   firstBuyAward,
			Hongbao:         honbao,
			Delflag:         0,
			StartTime:       start,
			EndTime:         end,
			CreateTime:      end,
			UpdateTime:      end,
		}
		agBill.Insert(db.Pool)
	}

	dao.UpdateAgPayStatus(db.Pool)
}
