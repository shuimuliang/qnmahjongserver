package hb

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
	Daifengpai           bool // 带风牌(勾选后牌组中包含东南西北风和中发白牌)
	Kechipai             bool // 可吃牌(上家打出牌，与自己手中的牌能组成顺子时可选择吃此张牌)
	Dianpaokehu          bool // 点炮可胡(可接别人打出的牌胡牌)
	Huangzhuanghuanggang bool // 黄庄黄杠(流局不计算杠分)

	// 加分
	Daizhuangxian bool // 带庄闲(庄家输赢结算时分数皆x2)
	Menqing       bool // 门清(胡牌时没有吃、碰、明杠过，得分x2)
	Biankadiao    bool // 边卡吊(胡边、卡、吊分数x2。边：手里有12或89牌胡37牌；卡：胡顺子中间那张牌；吊：手里有一张将牌，单吊另一张牌做将)
	Zhuowukui     bool // 捉五魁(手中有4万和6万胡卡5万得分x4)
	Daihuner      bool // 带混儿(混儿牌可以代替任何牌。抓完牌后，庄家翻出一张牌，此张牌的数字加一，为混儿牌)
	Suhu          bool // 素胡(胡牌时牌中没有混儿牌，得分x2)
	Hunerdiao     bool // 混儿吊(混儿做将牌，单吊任意一张牌胡牌，只能自摸胡，得分x2)

	// 算分
	Dianpaoyijiachu  bool // 点炮一家出(只有点炮者输分，其他两家不输分)
	Dianpaosanjiachu bool // 点炮三家出(有人点炮后，三家皆输分，点炮者多输一倍的分数)
	Dianpaodabao     bool // 点炮大包(点炮者输三倍的分数，其他两家不输分)
}

var configIdx = map[string]int{
	"Daifengpai":           0,
	"Kechipai":             1,
	"Dianpaokehu":          2,
	"Huangzhuanghuanggang": 3,
	"Daizhuangxian":        4,
	"Menqing":              5,
	"Biankadiao":           6,
	"Zhuowukui":            7,
	"Daihuner":             8,
	"Suhu":                 9,
	"Hunerdiao":            10,
	"Dianpaoyijiachu":      11,
	"Dianpaosanjiachu":     12,
	"Dianpaodabao":         13,
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

	var bitmap [14]bool
	for _, idx := range configs {
		bitmap[idx-1] = true
	}
	config = Config{
		Jushu:                round,
		Fangka:               0,
		Youxibi:              0,
		Daifengpai:           bitmap[configIdx["Daifengpai"]],
		Kechipai:             bitmap[configIdx["Kechipai"]],
		Dianpaokehu:          bitmap[configIdx["Dianpaokehu"]],
		Huangzhuanghuanggang: bitmap[configIdx["Huangzhuanghuanggang"]],
		Daizhuangxian:        bitmap[configIdx["Daizhuangxian"]],
		Menqing:              bitmap[configIdx["Menqing"]],
		Biankadiao:           bitmap[configIdx["Biankadiao"]],
		Zhuowukui:            bitmap[configIdx["Zhuowukui"]],
		Daihuner:             bitmap[configIdx["Daihuner"]],
		Suhu:                 bitmap[configIdx["Suhu"]],
		Hunerdiao:            bitmap[configIdx["Hunerdiao"]],
		Dianpaoyijiachu:      bitmap[configIdx["Dianpaoyijiachu"]],
		Dianpaosanjiachu:     bitmap[configIdx["Dianpaosanjiachu"]],
		Dianpaodabao:         bitmap[configIdx["Dianpaodabao"]],
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

	// 点炮可胡检查
	if config.Dianpaokehu {
		trueCount := 0
		if config.Dianpaoyijiachu {
			trueCount++
		}
		if config.Dianpaosanjiachu {
			trueCount++
		}
		if config.Dianpaodabao {
			trueCount++
		}
		if trueCount != 1 {
			ok = false
			return
		}
	} else if config.Dianpaoyijiachu || config.Dianpaosanjiachu || config.Dianpaodabao {
		ok = false
		return
	}

	// 带混儿检查
	if !config.Daihuner {
		if config.Suhu || config.Hunerdiao {
			ok = false
			return
		}
	}
	return
}
