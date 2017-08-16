package kf

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
	Dianpaokehu   bool // 点炮可胡
	Zhinengzimohu bool // 只能自模胡
	Maipao        bool // 买跑

	// 加分
	Zimofanbei         bool // 自摸翻倍
	Angangfanbei       bool // 暗杠翻倍
	Gangshanghuajiabei bool // 杠上花加倍
	Qiduifanbei        bool // 七对翻倍
	Gangpao            bool // 杠跑
}

var configIdx = map[string]int{
	"Dianpaokehu":        0,
	"Zhinengzimohu":      1,
	"Maipao":             2,
	"Zimofanbei":         3,
	"Angangfanbei":       4,
	"Gangshanghuajiabei": 5,
	"Qiduifanbei":        6,
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
		Dianpaokehu:        bitmap[configIdx["Dianpaokehu"]],
		Zhinengzimohu:      bitmap[configIdx["Zhinengzimohu"]],
		Maipao:             bitmap[configIdx["Maipao"]],
		Zimofanbei:         bitmap[configIdx["Zimofanbei"]],
		Angangfanbei:       bitmap[configIdx["Angangfanbei"]],
		Gangshanghuajiabei: bitmap[configIdx["Gangshanghuajiabei"]],
		Qiduifanbei:        bitmap[configIdx["Qiduifanbei"]],
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

	// 胡牌检查
	if config.Dianpaokehu && config.Zhinengzimohu {
		ok = false
		return
	}

	if !config.Dianpaokehu && !config.Zhinengzimohu {
		ok = false
		return
	}

	// 杠跑检查
	if !config.Maipao && config.Gangpao {
		ok = false
		return
	}
	return
}
