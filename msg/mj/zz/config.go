package zz

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"qnmahjong/pf"
)

type Config struct {
	// 局数
	Jushu int32 // 8局／16局

	// 房卡
	Fangka int32 // 3张／6张

	// 游戏币
	Youxibi int32 // 30个／60个

	// 玩法
	Daifengpai  bool // 带风牌
	Daihuner    bool // 带混儿
	Maipao      bool // 买跑
	Dianpaokehu bool // 点炮可胡

	// 加分
	Qiduijiabei        bool // 七对加倍
	Gangshanghuajiabei bool // 杠上花加倍
	Zhuangjiajiadi     bool // 庄家加底
	Gangpao            bool // 杠跑
}

var configIdx = map[string]int{
	"Daifengpai":         0,
	"Daihuner":           1,
	"Maipao":             2,
	"Dianpaokehu":        3,
	"Qiduijiabei":        4,
	"Gangshanghuajiabei": 5,
	"Zhuangjiajiadi":     6,
	"Gangpao":            7,
}

func checkConfig(id int32, send *pf.EnterRoomSend) (config Config, ok bool) {
	configs := send.Configs
	mjType := send.MjType
	round := send.Round
	costType := send.CostType

	cost, ok := cache.GetCosts(mjType, round, costType)
	if !ok {
		return
	}

	var bitmap [8]bool
	for _, idx := range configs {
		bitmap[idx-1] = true
	}
	config = Config{
		Jushu:              round,
		Fangka:             0,
		Youxibi:            0,
		Daifengpai:         bitmap[configIdx["Daifengpai"]],
		Daihuner:           bitmap[configIdx["Daihuner"]],
		Maipao:             bitmap[configIdx["Maipao"]],
		Dianpaokehu:        bitmap[configIdx["Dianpaokehu"]],
		Qiduijiabei:        bitmap[configIdx["Qiduijiabei"]],
		Gangshanghuajiabei: bitmap[configIdx["Gangshanghuajiabei"]],
		Zhuangjiajiadi:     bitmap[configIdx["Zhuangjiajiadi"]],
		Gangpao:            bitmap[configIdx["Gangpao"]],
	}

	// 房卡金币检查
	if costType == def.CostCoin {
		config.Youxibi = cost
		if !cache.CheckCoins(id, cost) {
			ok = false
			return
		}
	}

	if costType == def.CostCard {
		config.Fangka = cost
		if !cache.CheckCards(id, cost) {
			ok = false
			return
		}
	}

	// 杠跑检查
	if !config.Maipao && config.Gangpao {
		ok = false
		return
	}
	return
}
