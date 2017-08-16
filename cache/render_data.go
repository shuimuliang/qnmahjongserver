package cache

import (
	"fmt"
	"qnmahjong/db/dao"
)

type UserCardsRD struct {
	StartDate   string
	EndDate     string
	CardsCnt    int32
	MoneyCnt    int32
	OrderLength int32
	OrderList   []*dao.AgPay
}

func (this *UserCardsRD) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

type UserBalanceRD struct {
	Balance1 int32 // 今日提成
	Balance2 int32 // 本周总提成
	Balance3 int32 // 本周下级代理提成
	Balance4 int32 // 本周我直接销售的玉提成
	Balance5 int32 // 本周活动首充赠送
	Balance6 int32 // 500元分销红包
	Balance7 int32 // 总打款金额
}

func (this *UserBalanceRD) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

type QueryOrderRD struct {
	StartDate   string
	EndDate     string
	PlayerID    int32
	CardsCnt    int32
	MoneyCnt    int32
	OrderLength int32
	OrderList   []*dao.AgPay
}

func (this *QueryOrderRD) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

type SaleIndexRD struct {
	*AgAgent
	CardCnt    int32
	BalanceCnt int32
	AgentCnt   int32
	PlayerCnt  int32
}

func (this *SaleIndexRD) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

type VerifyBalanceRD struct {
	Start string
	End   string
}

type VerifyBalanceDetailRD struct {
	*AgAgent
	*dao.AgBill
}

func (this *VerifyBalanceDetailRD) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}
