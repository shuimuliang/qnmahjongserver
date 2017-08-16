package hb

type Fanshu struct {
	// 加分牌型
	Qidui                   bool // 七对x2
	Pengpenghu              bool // 碰碰胡x2
	Qingyise                bool // 清一色x2
	Qingfeng                bool // 清风x2
	Yitiaolong              bool // 一条龙x2
	Haohuaqidui             bool // 豪华七对x4（不再算七对）
	Chaojihaohuaqidui       bool // 超级豪华七对x8（不再算七对，豪华七对）
	Zhizunchaojihaohuaqidui bool // 至尊超级豪华七对x16（不再算七对，豪华七对，超级豪华七对）
	Shisanyao               bool // 十三幺x10
	// 额外加分
	Menqing         bool // 门清x2
	Bian            bool // 边x2
	Ka              bool // 卡x2
	Diao            bool // 吊x2
	Gangshangkaihua bool // 杠上开花x2
	Haidilaoyue     bool // 海底捞月x2
	Zhuowukui       bool // 捉五魁x4（不再算卡）
	Suhu            bool // 素胡x2
	Hunerdiao       bool // 混儿吊x2（不再算吊）
	Daizhuangxian   bool // 带庄闲，庄家输赢x2
	Zimo            bool // 自摸x2
}

func (tiles Tiles) GetFanshu7Pair() (fanshu Fanshu) {
	fanshu = Fanshu{}
	var count int
	for i := 0; i < len(tiles)-1; i++ {
		if tiles[i].Unicode != tiles[i+1].Unicode {
			count++
		}
	}
	if count == 3 {
		fanshu.Zhizunchaojihaohuaqidui = true
	} else if count == 4 {
		fanshu.Chaojihaohuaqidui = true
	} else if count == 5 {
		fanshu.Haohuaqidui = true
	} else {
		fanshu.Qidui = true
	}

	if tiles.IsQingyise() {
		fanshu.Qingyise = true
	}

	if tiles.IsQingfeng() {
		fanshu.Qingfeng = true
	}
	return
}

func (tiles Tiles) GetFanshu13yao() (fanshu Fanshu) {
	fanshu = Fanshu{}
	fanshu.Shisanyao = true
	return
}

func (tiles Tiles) GetFanshuElse(config *MJConfig) (fanshu Fanshu) {
	fanshu = Fanshu{}
	if tiles.IsPengPengHu() {
		fanshu.Pengpenghu = true
	}
	if tiles.IsQingyise() {
		fanshu.Qingyise = true
	}
	if tiles.IsQingfeng() {
		fanshu.Qingfeng = true
	}
	if tiles.IsYitiaolong() {
		fanshu.Yitiaolong = true
	}
	if config.Config.Biankadiao {
		if tiles.IsBian() {
			fanshu.Bian = true
		}
		if tiles.IsKa() {
			fanshu.Ka = true
		}
		if tiles.IsDiao() {
			fanshu.Diao = true
		}
	}
	if config.Config.Zhuowukui {
		if tiles.IsZhuoWuKui() {
			fanshu.Ka = false
			fanshu.Zhuowukui = true
		}
	}
	return
}

func (tiles Tiles) GetFanshuElseWithHun(config *MJConfig) (fanshu Fanshu) {
	fanshu = Fanshu{}
	if tiles.IsPengPengHu() {
		fanshu.Pengpenghu = true
	}
	if tiles.IsQingyise() {
		fanshu.Qingyise = true
	}
	if tiles.IsQingfeng() {
		fanshu.Qingfeng = true
	}
	if tiles.IsYitiaolong() {
		fanshu.Yitiaolong = true
	}
	// if config.Config.Biankadiao {
	// 	if tiles.IsBian(huTile) {
	// 		fanshu.Bian = true
	// 	}
	// 	if tiles.IsKa(huTile) {
	// 		fanshu.Ka = true
	// 	}
	// 	if tiles.IsDiao(huTile) {
	// 		fanshu.Diao = true
	// 	}
	// }
	// if config.Config.Zhuowukui {
	// 	if tiles.IsZhuoWuKui(huTile) {
	// 		fanshu.Zhuowukui = true
	// 		fanshu.Ka = false
	// 	}
	// }
	// if config.Config.Hunerdiao {
	// 	if tiles.IsHunerdiao(huTile) {
	// 		fanshu.Hunerdiao = true
	// 		fanshu.Diao = false
	// 	}
	// }
	return
}

func (tiles Tiles) IsPengPengHu() (ok bool) {
	var count int
	for i := 0; i < len(tiles)-4; i = i + 3 {
		if tiles[i].Unicode == tiles[i+1].Unicode && tiles[i].Unicode == tiles[i+2].Unicode {
			count++
		}
	}
	if count == 4 {
		ok = true
	}
	return
}

func (tiles Tiles) IsQingyise() (ok bool) {
	var count int
	if tiles[0].Type == Characters || tiles[0].Type == Bamboos || tiles[0].Type == Circles {
		for i := 1; i < len(tiles); i++ {
			if tiles[i].Type == tiles[0].Type {
				count++
			}
		}
	}
	if count == 13 {
		ok = true
	}
	return
}

func (tiles Tiles) IsQingfeng() (ok bool) {
	var count int
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Type == Wind {
			count++
		}
	}
	if count == 14 {
		ok = true
	}
	return
}

func (tiles Tiles) IsYitiaolong() (ok bool) {
	var wan = []Tile{
		{Characters, One, OneOfCharacters, 1, false, false},
		{Characters, Two, TwoOfCharacters, 2, false, false},
		{Characters, Three, ThreeOfCharacters, 3, false, false},
		{Characters, Four, FourOfCharacters, 4, false, false},
		{Characters, Five, FiveOfCharacters, 5, false, false},
		{Characters, Six, SixOfCharacters, 6, false, false},
		{Characters, Seven, SevenOfCharacters, 7, false, false},
		{Characters, Eight, EightOfCharacters, 8, false, false},
		{Characters, Nine, NineOfCharacters, 9, false, false},
	}

	var bing = []Tile{
		{Circles, One, OneOfCircles, 11, false, false},
		{Circles, Two, TwoOfCircles, 12, false, false},
		{Circles, Three, ThreeOfCircles, 13, false, false},
		{Circles, Four, FourOfCircles, 14, false, false},
		{Circles, Five, FiveOfCircles, 15, false, false},
		{Circles, Six, SixOfCircles, 16, false, false},
		{Circles, Seven, SevenOfCircles, 17, false, false},
		{Circles, Eight, EightOfCircles, 18, false, false},
		{Circles, Nine, NineOfCircles, 19, false, false},
	}

	var tiao = []Tile{
		{Bamboos, One, OneOfBamboos, 21, false, false},
		{Bamboos, Two, TwoOfBamboos, 22, false, false},
		{Bamboos, Three, ThreeOfBamboos, 23, false, false},
		{Bamboos, Four, FourOfBamboos, 24, false, false},
		{Bamboos, Five, FiveOfBamboos, 25, false, false},
		{Bamboos, Six, SixOfBamboos, 26, false, false},
		{Bamboos, Seven, SevenOfBamboos, 27, false, false},
		{Bamboos, Eight, EightOfBamboos, 28, false, false},
		{Bamboos, Nine, NineOfBamboos, 29, false, false},
	}

	var i, j, k int
	for _, tile := range tiles {
		if tile.Unicode == wan[i].Unicode {
			i++
			if i == 9 {
				break
			}
		}
		if tile.Unicode == bing[j].Unicode {
			j++
			if j == 9 {
				break
			}
		}
		if tile.Unicode == tiao[k].Unicode {
			k++
			if k == 9 {
				break
			}
		}
	}
	if i == 9 || j == 9 || k == 9 {
		ok = true
	}
	return
}

func (seat *MJSeat) IsMenqing() (ok bool) {
	if len(seat.AnGangTiles)+len(seat.MingGangTiles)+len(seat.BuGangTiles)+len(seat.ChiTiles)+len(seat.PengTiles) == 0 {
		ok = true
	}
	return
}

func (tiles Tiles) IsBian() (ok bool) {
	for i := 0; i < len(tiles)-4; i = i + 3 {
		if tiles[i+2].IsHu && tiles[i+2].Value == 3 {
			if tiles[i].Unicode == tiles[i+2].Unicode-2 &&
				tiles[i+1].Unicode == tiles[i+2].Unicode-1 {
				ok = true
				break
			}
		}
	}
	for i := 0; i < len(tiles)-4; i = i + 3 {
		if tiles[i].IsHu && tiles[i].Value == 7 {
			if tiles[i+1].Unicode == tiles[i].Unicode+1 &&
				tiles[i+2].Unicode == tiles[i].Unicode+2 {
				ok = true
				break
			}
		}
	}
	return
}

func (tiles Tiles) IsKa() (ok bool) {
	for i := 0; i < len(tiles)-4; i = i + 3 {
		if tiles[i+1].IsHu {
			if tiles[i].Unicode == tiles[i+1].Unicode-1 &&
				tiles[i+2].Unicode == tiles[i+1].Unicode+1 {
				ok = true
				break
			}
		}
	}
	return
}

func (tiles Tiles) IsDiao() (ok bool) {
	l := len(tiles)
	if tiles[l-1].IsHu || tiles[l-2].IsHu {
		ok = true
	}
	return
}

func (tiles Tiles) IsHunerdiao(tile Tile) (ok bool) {
	if tiles[len(tiles)-1].Unicode == tile.Unicode {
		if tiles[len(tiles)-1].IsHun || tiles[len(tiles)-2].IsHun {
			ok = true
		}
	}
	return
}

func (tiles Tiles) IsZhuoWuKui() (ok bool) {
	for i := 0; i < len(tiles)-4; i = i + 3 {
		if tiles[i+1].IsHu {
			if tiles[i+1].Unicode == FiveOfCharacters {
				if tiles[i].Unicode == tiles[i+1].Unicode-1 &&
					tiles[i+2].Unicode == tiles[i+1].Unicode+1 {
					ok = true
					break
				}
			}
		}
	}
	return
}

func (fanshu Fanshu) GetScore() (score int32, huTypes []int32) {
	score = 1
	if fanshu.Qidui {
		score = score * 2
		huTypes = append(huTypes, 1)
	}
	if fanshu.Pengpenghu {
		score = score * 2
		huTypes = append(huTypes, 2)
	}
	if fanshu.Qingyise {
		score = score * 2
		huTypes = append(huTypes, 3)
	}
	if fanshu.Qingfeng {
		score = score * 2
		huTypes = append(huTypes, 4)
	}
	if fanshu.Yitiaolong {
		score = score * 2
		huTypes = append(huTypes, 5)
	}
	if fanshu.Haohuaqidui {
		score = score * 4
		huTypes = append(huTypes, 6)
	}
	if fanshu.Chaojihaohuaqidui {
		score = score * 8
		huTypes = append(huTypes, 7)
	}
	if fanshu.Zhizunchaojihaohuaqidui {
		score = score * 16
		huTypes = append(huTypes, 8)
	}
	if fanshu.Shisanyao {
		score = score * 10
		huTypes = append(huTypes, 9)
	}
	if fanshu.Menqing {
		score = score * 2
		huTypes = append(huTypes, 10)
	}
	if fanshu.Bian {
		score = score * 2
		huTypes = append(huTypes, 11)
	}
	if fanshu.Ka {
		score = score * 2
		huTypes = append(huTypes, 12)
	}
	if fanshu.Diao {
		score = score * 2
		huTypes = append(huTypes, 13)
	}
	if fanshu.Gangshangkaihua {
		score = score * 2
		huTypes = append(huTypes, 14)
	}
	if fanshu.Haidilaoyue {
		score = score * 2
		huTypes = append(huTypes, 15)
	}
	if fanshu.Zhuowukui {
		score = score * 4
		huTypes = append(huTypes, 16)
	}
	if fanshu.Suhu {
		score = score * 2
		huTypes = append(huTypes, 17)
	}
	if fanshu.Hunerdiao {
		score = score * 2
		huTypes = append(huTypes, 18)
	}
	if fanshu.Daizhuangxian {
		// score=score * 2
		huTypes = append(huTypes, 19)
	}
	if fanshu.Zimo {
		score = score * 2
		huTypes = append(huTypes, 20)
	}
	return
}
