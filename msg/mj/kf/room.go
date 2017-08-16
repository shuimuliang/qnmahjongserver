package kf

import (
	"fmt"
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/redis"
	"qnmahjong/util"
	"math/rand"
	"sync"
	"time"
)

// IsZiMo 自摸
func (s *MJSeat) IsZiMo(config *MJConfig) (huTiles Tiles, fanshu Fanshu, ok bool) {
	tiles := make(Tiles, len(s.HandTiles))
	copy(tiles, s.HandTiles)
	zimoTile := Tile{
		Type:    s.DrawTile.Type,
		Value:   s.DrawTile.Value,
		Unicode: s.DrawTile.Unicode,
		Front:   s.DrawTile.Front,
		IsHu:    true,
	}
	tiles = append(tiles, zimoTile)
	huTiles, fanshu, ok = tiles.IsHupai(s, config)
	return
}

// IsDianPao 点炮
func (s *MJSeat) IsDianPao(config *MJConfig, dianpaoTile *Tile) (huTiles Tiles, fanshu Fanshu, ok bool) {
	tiles := make(Tiles, len(s.HandTiles))
	copy(tiles, s.HandTiles)
	dianpaoTile.IsHu = true
	tiles = append(tiles, *dianpaoTile)
	huTiles, fanshu, ok = tiles.IsHupai(s, config)
	return
}

// IsAnGang 暗杠
func (s *MJSeat) IsAnGang() (anGangTiles []Tile) {
	tiles := make(Tiles, len(s.HandTiles))
	copy(tiles, s.HandTiles)
	tiles = append(tiles, *s.DrawTile)
	tiles.SortTiles()
	if len(tiles) < 4 {
		return
	}

	for i := 0; i <= len(tiles)-4; i++ {
		if tiles[i].Unicode == tiles[i+1].Unicode &&
			tiles[i].Unicode == tiles[i+2].Unicode &&
			tiles[i].Unicode == tiles[i+3].Unicode {
			anGangTiles = append(anGangTiles, tiles[i])
		}
	}
	return
}

// IsBuGang 补杠
func (s *MJSeat) IsBuGang() (buGangTiles []Tile) {
	tiles := make(Tiles, len(s.HandTiles))
	copy(tiles, s.HandTiles)
	tiles = append(tiles, *s.DrawTile)
	for i := 0; i < len(s.PengTiles); i = i + 3 {
		for _, tile := range tiles {
			if tile.Unicode == s.PengTiles[i].Unicode {
				buGangTiles = append(buGangTiles, tile)
				break
			}
		}
	}
	return
}

// IsMingGang 明杠
func (s *MJSeat) IsMingGang(tile Tile) (ok bool) {
	tiles := s.HandTiles
	if len(tiles) < 3 {
		return
	}

	var count = 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Unicode == tile.Unicode {
			count++
		}
	}
	if count >= 3 {
		ok = true
	}
	return
}

// IsPeng 碰
func (s *MJSeat) IsPeng(tile Tile) (ok bool) {
	tiles := s.HandTiles
	if len(tiles) < 2 {
		return
	}

	var count = 0
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Unicode == tile.Unicode {
			count++
		}
	}
	if count >= 2 {
		ok = true
	}
	return
}

func (s *MJSeat) GetCardsInfo() (cardsInfo *pf.CardsInfo) {
	mingGangCards := s.MingGangTiles.GetFrontCards()
	mingGangCards = append(mingGangCards, s.BuGangTiles.GetFrontCards()...)
	front := int32(-1)
	if s.DrawTile != nil {
		front = s.DrawTile.Front
	}
	cardsInfo = &pf.CardsInfo{
		Pos:           s.Pos + 1,
		PengCards:     s.PengTiles.GetFrontCards(),
		ChiCards:      s.ChiTiles.GetFrontCards(),
		AnGangCards:   s.AnGangTiles.GetFrontCards(),
		MingGangCards: mingGangCards,
		DiscardCards:  s.DeskTiles.GetFrontCards(),
		LeftCards:     s.HandTiles.GetFrontCards(),
		DrawCard:      front,
		PaoScore:      s.PaoScore,
	}
	if !s.IsBuyPao {
		cardsInfo.PaoScore = -1
	}
	return
}

func (s *MJSeat) GetGangCards() (cards []int32) {
	cards = append(cards, s.AnGangTiles.GetFrontCards()...)
	cards = append(cards, s.MingGangTiles.GetFrontCards()...)
	cards = append(cards, s.BuGangTiles.GetFrontCards()...)
	return
}

func (s *MJSeat) GetRulerTiles() (tiles Tiles) {
	for i := 0; i < len(s.AnGangTiles); i = i + 4 {
		tiles = append(tiles, s.AnGangTiles[i], s.AnGangTiles[i], s.AnGangTiles[i])
	}
	for i := 0; i < len(s.MingGangTiles); i = i + 4 {
		tiles = append(tiles, s.MingGangTiles[i], s.MingGangTiles[i], s.MingGangTiles[i])
	}
	for i := 0; i < len(s.BuGangTiles); i = i + 4 {
		tiles = append(tiles, s.BuGangTiles[i], s.BuGangTiles[i], s.BuGangTiles[i])
	}
	tiles = append(tiles, s.PengTiles...)
	tiles = append(tiles, s.ChiTiles...)
	return
}

func (s *MJSeat) GetSelfHupaiList(room *MJRoom, afterGang bool) (operationList []*pf.OperationItem) {
	if huTiles, fanshu, ok := s.IsZiMo(room.RoomConfig); ok {
		// 自摸
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.ZimoHu,
			Keycard:   s.DrawTile.Front,
			CardsList: nil,
			SrcPos:    s.Pos + 1,
			DesPos:    s.Pos + 1,
		})
		fanshu.Gangshanghuajiabei = afterGang
		if room.RoomConfig.Config.Zimofanbei {
			fanshu.Zimofanbei = true
		}
		room.RoomHupai.IsZiMo = true
		room.RoomHupai.IsZiMoHupai = false
		room.RoomHupai.ZiMoPos = s.Pos
		room.RoomHupai.ZiMoTile = s.DrawTile
		room.RoomHupai.ZiMoHuTiles = huTiles
		room.RoomHupai.ZiMoHuFanshu = fanshu

		// 过
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.Pass,
			Keycard:   s.DrawTile.Front,
			CardsList: nil,
			SrcPos:    s.Pos + 1,
			DesPos:    s.Pos + 1,
		})
	}
	return
}

func (s *MJSeat) GetSelfGangList() (operationList []*pf.OperationItem) {
	// 暗杠
	anGangtiles := s.IsAnGang()
	for _, tile := range anGangtiles {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.AnGang,
			Keycard:   tile.Front,
			CardsList: nil,
			SrcPos:    s.Pos + 1,
			DesPos:    s.Pos + 1,
		})
	}
	// 补杠
	buGangtiles := s.IsBuGang()
	for _, tile := range buGangtiles {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.BuGang,
			Keycard:   tile.Front,
			CardsList: nil,
			SrcPos:    s.Pos + 1,
			DesPos:    s.Pos + 1,
		})
	}
	// 过
	if len(operationList) > 0 {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.Pass,
			Keycard:   s.DrawTile.Front,
			CardsList: nil,
			SrcPos:    s.Pos + 1,
			DesPos:    s.Pos + 1,
		})
	}
	return
}

func (s *MJSeat) GetOtherHupaiList(room *MJRoom, tile Tile, pos int32) (operationList []*pf.OperationItem) {
	// 点炮
	if room.RoomConfig.Config.Dianpaokehu {
		if huTiles, fanshu, ok := s.IsDianPao(room.RoomConfig, &tile); ok {
			operationList = append(operationList, &pf.OperationItem{
				Type:      def.DianpaoHu,
				Keycard:   tile.Front,
				CardsList: nil,
				SrcPos:    pos + 1,
				DesPos:    s.Pos + 1,
			})
			room.RoomHupai.IsDianPao = true
			room.RoomHupai.IsDianPaoHupai = false
			room.RoomHupai.DianPaoPos = pos
			room.RoomHupai.JiePaoPos = s.Pos
			room.RoomHupai.DianPaoTile = &tile
			room.RoomHupai.DianPaoHuTiles = huTiles
			room.RoomHupai.DianPaoHuFanshu = fanshu
			// 过
			operationList = append(operationList, &pf.OperationItem{
				Type:      def.Pass,
				Keycard:   tile.Front,
				CardsList: nil,
				SrcPos:    pos + 1,
				DesPos:    s.Pos + 1,
			})
		}
	}
	return
}

func (s *MJSeat) GetOtherGangPengList(tile Tile, pos int32) (operationList []*pf.OperationItem) {
	// 明杠
	ok := s.IsMingGang(tile)
	if ok {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.MingGang,
			Keycard:   tile.Front,
			CardsList: nil,
			SrcPos:    pos + 1,
			DesPos:    s.Pos + 1,
		})
	}

	// 碰
	ok = s.IsPeng(tile)
	if ok {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.Peng,
			Keycard:   tile.Front,
			CardsList: nil,
			SrcPos:    pos + 1,
			DesPos:    s.Pos + 1,
		})
	}

	// 过
	if len(operationList) > 0 {
		operationList = append(operationList, &pf.OperationItem{
			Type:      def.Pass,
			Keycard:   tile.Front,
			CardsList: nil,
			SrcPos:    pos + 1,
			DesPos:    s.Pos + 1,
		})
	}
	return
}

func (s *MJSeat) AddGangTypes(angangfanbei bool, huTypes *[]int32) {
	if len(s.IsMingGangTiles) > 0 {
		*huTypes = append(*huTypes, 207)
	}
	if len(s.MingGangTiles) > 0 {
		*huTypes = append(*huTypes, 208)
	}
	if len(s.AnGangTiles) > 0 {
		if angangfanbei {
			*huTypes = append(*huTypes, 214)
		} else {
			*huTypes = append(*huTypes, 209)
		}
	}
	if len(s.IsBuGangTiles) > 0 {
		*huTypes = append(*huTypes, 210)
	}
	if len(s.BuGangTiles) > 0 {
		*huTypes = append(*huTypes, 211)
	}
}

func (t *MJTiles) ShuffleTiles() {
	num := len(t.Tiles)
	result := make(Tiles, num)
	copy(result, t.Tiles)
	for i := 0; i < num; i++ {
		r := i + rand.Intn(num-i)
		result[r], result[i] = result[i], result[r]
	}
	t.Tiles = result

	// 测试使用
	// t.Tiles = TilesForTest

	t.DrawPos = 0
	t.LeftCnt = int32(len(t.Tiles))
	t.GangCnt = 0
}

func (t *MJTiles) DealTiles() (result Tiles) {
	result = make([]Tile, def.DiscardCount)
	newPos := t.DrawPos + def.DiscardCount
	copy(result, t.Tiles[t.DrawPos:newPos])
	t.DrawPos = newPos
	t.LeftCnt -= def.DiscardCount
	return result
}

func (t *MJTiles) DealOneTile() (result *Tile) {
	result = &t.Tiles[t.DrawPos]
	t.DrawPos++
	t.LeftCnt--
	return
}

func (t *MJTiles) IsLiuju() (ok bool) {
	ok = t.LeftCnt <= 0
	return
}

func (t *MJStatus) CheckOperation(operation *pf.OperationItem) (ok bool) {
	for _, item := range t.operationNoticeRecv.OperationList {
		if item.Equal(operation) {
			ok = true
			break
		}
	}
	return
}

func (t *MJStatus) GetOperation() (opType int32) {
	for _, item := range t.operationNoticeRecv.OperationList {
		opType = item.GetType()
		break
	}
	return
}

func (r *MJRoom) InRoom(playerID, roomID int32, otherIDs *[]int32) (pos int32, ok bool) {
	for _, seat := range r.RoomSeats {
		if seat == nil {
			continue
		}
		if seat.ID == playerID {
			pos = seat.Pos
			ok = true
		} else if otherIDs != nil {
			*otherIDs = append(*otherIDs, seat.ID)
		}
	}
	return
}

func (r *MJRoom) SuanFen() {
	if r.RoomHupai.IsZiMoHupai {
		ziMoPos := r.RoomHupai.ZiMoPos
		ziMoSeat := r.RoomSeats[ziMoPos]
		fanshu := r.RoomHupai.ZiMoHuFanshu
		for _, seat := range r.RoomSeats {
			if seat.Pos != ziMoPos {
				paoScore := ziMoSeat.PaoScore + seat.PaoScore
				huScore := 1 + paoScore

				if fanshu.Qiduifanbei {
					huScore *= 2
				}
				if fanshu.Gangshanghuajiabei {
					huScore *= 2
				}
				if fanshu.Zimofanbei {
					huScore *= 2
				}
				ziMoSeat.HuScore += huScore
				seat.HuScore -= huScore
			}
		}
	} else if r.RoomHupai.IsDianPaoHupai {
		dianPaoPos := r.RoomHupai.DianPaoPos
		JiePaoPos := r.RoomHupai.JiePaoPos
		dianPaoSeat := r.RoomSeats[dianPaoPos]
		jiePaoSeat := r.RoomSeats[JiePaoPos]
		paoScore := dianPaoSeat.PaoScore + jiePaoSeat.PaoScore

		huScore := 1 + paoScore

		fanshu := r.RoomHupai.DianPaoHuFanshu
		if fanshu.Gangshanghuajiabei {
			huScore *= 2
		}

		dianPaoSeat.HuScore -= huScore
		jiePaoSeat.HuScore += huScore
	} else if r.RoomHupai.IsMultiDianPaoHupai {
		for _, jiePaoPos := range r.RoomHupai.MultiJiePaoHupaiPos {
			dianPaoPos := r.RoomHupai.MultiDianPaoPos
			dianPaoSeat := r.RoomSeats[dianPaoPos]
			jiePaoSeat := r.RoomSeats[jiePaoPos]
			paoScore := dianPaoSeat.PaoScore + jiePaoSeat.PaoScore

			huScore := 1 + paoScore

			for i, pos := range r.RoomHupai.MultiJiePaoPos {
				if pos == jiePaoPos {
					fanshu := r.RoomHupai.MultiDianPaoHuFanshu[i]
					if fanshu.Gangshanghuajiabei {
						huScore *= 2
					}

					dianPaoSeat.HuScore -= huScore
					jiePaoSeat.HuScore += huScore
					break
				}
			}
		}
	}
	return
}

func (r *MJRoom) GetResultItems() (resultItems []*pf.ResultItem) {
	for _, seat := range r.RoomSeats {
		drawCard := int32(0)
		isHu := int32(0)
		huTypes := []int32{}
		if r.RoomHupai.IsZiMoHupai {
			if seat.Pos == r.RoomHupai.ZiMoPos {
				isHu = 1
				huTypes = r.RoomHupai.ZiMoHuFanshu.GetHuTypes()
				if r.RoomConfig.Config.Zimofanbei {
					huTypes = append(huTypes, 213)
				} else {
					huTypes = append(huTypes, 203)
				}
				seat.ZiMoCnt++
				drawCard = r.RoomHupai.ZiMoTile.Front
				seat.HandTiles = r.RoomHupai.ZiMoHuTiles.DelHuCard()
			}
			seat.AddGangTypes(r.RoomConfig.Config.Angangfanbei, &huTypes)
		} else if r.RoomHupai.IsDianPaoHupai {
			if seat.Pos == r.RoomHupai.JiePaoPos {
				isHu = 1
				huTypes = r.RoomHupai.DianPaoHuFanshu.GetHuTypes()
				huTypes = append(huTypes, 202)
				seat.JiePaoCnt++
				drawCard = r.RoomHupai.DianPaoTile.Front
				seat.HandTiles = r.RoomHupai.DianPaoHuTiles.DelHuCard()
			} else if seat.Pos == r.RoomHupai.DianPaoPos {
				seat.DiaoPaoCnt++
				huTypes = append(huTypes, 201)
			}
			seat.AddGangTypes(r.RoomConfig.Config.Angangfanbei, &huTypes)
		} else if r.RoomHupai.IsMultiDianPaoHupai {
			isMultiHu := false
			index := int32(0)
			for _, jiePaoPos := range r.RoomHupai.MultiJiePaoHupaiPos {
				if seat.Pos == jiePaoPos {
					isMultiHu = true
					for i, pos := range r.RoomHupai.MultiJiePaoPos {
						if pos == jiePaoPos {
							index = int32(i)
							break
						}
					}
					break
				}
			}

			if isMultiHu {
				isHu = 1
				huTypes = r.RoomHupai.MultiDianPaoHuFanshu[index].GetHuTypes()
				huTypes = append(huTypes, 202)
				seat.JiePaoCnt++
				drawCard = r.RoomHupai.MultiDianPaoTile.Front
				seat.HandTiles = r.RoomHupai.MultiDianPaoHuTiles[index].DelHuCard()
			} else if seat.Pos == r.RoomHupai.MultiDianPaoPos {
				seat.DiaoPaoCnt++
				huTypes = append(huTypes, 201)
			}
			seat.AddGangTypes(r.RoomConfig.Config.Angangfanbei, &huTypes)
		} else {
			// 流局不算杠分
			seat.GangScore = 0
		}

		seat.Score += seat.HuScore + seat.GangScore
		resultItems = append(resultItems, &pf.ResultItem{
			Pos:       seat.Pos + 1,
			PengCards: seat.PengTiles.GetFrontCards(),
			ChiCards:  seat.ChiTiles.GetFrontCards(),
			GangCards: seat.GetGangCards(),
			LeftCards: seat.HandTiles.GetFrontCards(),
			DrawCard:  drawCard,
			Score:     seat.HuScore + seat.GangScore,
			IsHu:      isHu,
			HuTypes:   huTypes,
			HuScore:   seat.HuScore,
			GangScore: seat.GangScore,
			PaoScore:  seat.PaoScore,
			CurScore:  seat.Score,
		})

		// 战绩
		l := len(r.RoomRecord.RoundList)
		r.RoomRecord.RoundList[l-1].ScoreList[seat.Pos] = seat.HuScore + seat.GangScore
		switch seat.Pos + 1 {
		case East:
			r.DaoRecord.EastID = seat.ID
			r.DaoRecord.EastScore = seat.HuScore + seat.GangScore
		case South:
			r.DaoRecord.SouthID = seat.ID
			r.DaoRecord.SouthScore = seat.HuScore + seat.GangScore
		case West:
			r.DaoRecord.WestID = seat.ID
			r.DaoRecord.WestScore = seat.HuScore + seat.GangScore
		case North:
			r.DaoRecord.NorthID = seat.ID
			r.DaoRecord.NorthScore = seat.HuScore + seat.GangScore
		}
	}

	// 战绩记录
	err := r.DaoRecord.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "record", r.DaoRecord, 0, def.ErrInsertRecord)
	} else {
		r.RoomRecord.RoundList[r.RoomStatus.CurRound-1].RoundID = r.DaoRecord.RecordID
	}

	// 扣除房卡
	if !r.RoomStatus.IsSubCard {
		id := r.DaoRecord.CreateID
		cards := r.RoomConfig.Config.Fangka
		coins := r.RoomConfig.Config.Youxibi
		cache.HandleSubCard(id, cards, coins)
		r.RoomStatus.IsSubCard = true

		t := &dao.Treasure{
			PlayerID:   id,
			Reason:     def.CreateRoom,
			Coins:      coins,
			Cards:      cards,
			ChangeTime: time.Now(),
		}

		// 开房扣除
		err = t.Insert(db.Pool)
		if err != nil {
			util.LogError(err, "treasure", t, id, def.ErrInsertTreasure)
		}
	}
	return
}

func (r *MJRoom) GetSettleItems() (settleItems []*pf.SettleItem) {
	if r.DaoRecord.RecordID == 0 {
		for _, seat := range r.RoomSeats {
			switch seat.Pos + 1 {
			case East:
				r.DaoRecord.EastID = seat.ID
			case South:
				r.DaoRecord.SouthID = seat.ID
			case West:
				r.DaoRecord.WestID = seat.ID
			case North:
				r.DaoRecord.NorthID = seat.ID
			}
		}

		// 战绩记录
		err := r.DaoRecord.Insert(db.Pool)
		if err != nil {
			util.LogError(err, "record", r.DaoRecord, 0, def.ErrInsertRecord)
		} else {
			r.RoomRecord.RoundList[r.RoomStatus.CurRound-1].RoundID = r.DaoRecord.RecordID
		}
	}

	r.DaoRecord = &dao.Record{
		CreateID:   r.RoomSeats[0].ID,
		CreateTime: r.RoomConfig.CreateTime,
		RoomID:     r.RoomID,
		MjType:     r.RoomConfig.MJType,
		TotalRound: r.RoomConfig.Config.Jushu,
		StartTime:  time.Now(),
		CurRound:   0,
		EastID:     0,
		SouthID:    0,
		WestID:     0,
		NorthID:    0,
		EastScore:  0,
		SouthScore: 0,
		WestScore:  0,
		NorthScore: 0,
	}
	for _, seat := range r.RoomSeats {
		settleItems = append(settleItems, &pf.SettleItem{
			Pos:        seat.Pos + 1,
			Score:      seat.Score,
			ZiMoCnt:    seat.ZiMoCnt,
			JiePaoCnt:  seat.JiePaoCnt,
			DianPaoCnt: seat.DiaoPaoCnt,
		})
		// 战绩
		r.RoomRecord.ScoreList[seat.Pos] = seat.Score
		r.DaoRecord.CurRound = 0
		switch seat.Pos + 1 {
		case East:
			r.DaoRecord.EastID = seat.ID
			r.DaoRecord.EastScore = seat.Score
		case South:
			r.DaoRecord.SouthID = seat.ID
			r.DaoRecord.SouthScore = seat.Score
		case West:
			r.DaoRecord.WestID = seat.ID
			r.DaoRecord.WestScore = seat.Score
		case North:
			r.DaoRecord.NorthID = seat.ID
			r.DaoRecord.NorthScore = seat.Score
		}
	}
	// 战绩记录
	err := r.DaoRecord.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "record", r.DaoRecord, 0, def.ErrInsertRecord)
	}
	r.RecordToRedis()
	return
}

func (r *MJRoom) RecordToRedis() {
	for _, seat := range r.RoomSeats {
		redis.PutRecord(seat.ID, r.RoomRecord)
	}
}

func (r *MJRoom) AddAnGangScore(desPos int32) {
	gangSeat := r.RoomSeats[desPos]
	for _, seat := range r.RoomSeats {
		if seat.Pos != desPos {
			paoScore := gangSeat.PaoScore + seat.PaoScore
			if !r.RoomConfig.Config.Gangpao {
				paoScore = 0
			}
			if r.RoomConfig.Config.Angangfanbei {
				gangSeat.GangScore += 2 * (1 + paoScore)
				seat.GangScore -= 2 * (1 + paoScore)
			} else {
				gangSeat.GangScore += 1 + paoScore
				seat.GangScore -= 1 + paoScore
			}
		}
	}
}

func (r *MJRoom) AddMingGangScore(srcPos int32, desPos int32) {
	paoScore := r.RoomSeats[srcPos].PaoScore + r.RoomSeats[desPos].PaoScore
	if !r.RoomConfig.Config.Gangpao {
		paoScore = 0
	}
	r.RoomSeats[srcPos].GangScore -= 1 + paoScore
	r.RoomSeats[desPos].GangScore += 1 + paoScore
}

func (r *MJRoom) AddBuGangScore(srcPos int32, desPos int32) {
	paoScore := r.RoomSeats[srcPos].PaoScore + r.RoomSeats[desPos].PaoScore
	if !r.RoomConfig.Config.Gangpao {
		paoScore = 0
	}
	r.RoomSeats[srcPos].GangScore -= 1 + paoScore
	r.RoomSeats[desPos].GangScore += 1 + paoScore
}

func (r *MJRoom) ClearSeats() {
	for _, seat := range r.RoomSeats {
		seat.Prepared = false
		seat.IsBuyPao = false
		seat.PaoScore = 0

		seat.DrawTile = nil
		seat.DiscardTile = nil

		seat.HandTiles = nil
		seat.DeskTiles = nil

		seat.MingGangTiles = nil
		seat.AnGangTiles = nil
		seat.BuGangTiles = nil
		seat.PengTiles = nil
		seat.ChiTiles = nil

		seat.IsMingGangTiles = nil
		seat.IsBuGangTiles = nil
		seat.IsPengTiles = nil
		seat.IsChiTiles = nil

		seat.HuScore = 0
		seat.GangScore = 0
	}
}

func (r *MJRoom) ClearRoom() {
	// RoomStatus
	if r.RoomHupai.IsZiMoHupai {
		r.RoomStatus.BankerPos = r.RoomHupai.ZiMoPos
	}

	if r.RoomHupai.IsDianPaoHupai {
		r.RoomStatus.BankerPos = r.RoomHupai.JiePaoPos
	}

	if r.RoomHupai.IsMultiDianPaoHupai {
		for i := int32(1); i < def.RoomPlayerCount; i++ {
			pos := (r.RoomHupai.MultiDianPaoPos + i) % def.RoomPlayerCount
			isHu := false
			for _, huPos := range r.RoomHupai.MultiJiePaoHupaiPos {
				if huPos == pos {
					isHu = true
					break
				}
			}
			if isHu {
				r.RoomStatus.BankerPos = pos
				break
			}
		}
	}

	r.RoomStatus.CurPos = r.RoomStatus.BankerPos
	r.RoomStatus.CurRound++
	r.RoomStatus.LastPos = -2 // 前端-1为没有人出牌

	r.RoomStatus.operationNoticeRecv = nil
	r.RoomStatus.discardNoticerecv = nil
	r.RoomStatus.endRoundRecv = nil

	r.RoomStatus.GameStatus = def.GameStatusIsGaming
	r.RoomStatus.IsStart = true
	r.RoomStatus.IsGaming = true

	// RoomHupai
	r.RoomHupai = &MJHupai{}

	// RoomTiles
	r.RoomTiles.ShuffleTiles()

	// RoomSeats
	r.ClearSeats()

	// RoomRecord
	var curTime = time.Now()
	if r.RoomRecord == nil {
		var roomRecord = &pf.RecordRoom{
			RoomID:     r.RoomID,
			MjType:     r.RoomConfig.MJType,
			TotalRound: r.RoomConfig.Config.Jushu,
			CreateTime: r.RoomConfig.CreateTime.Format("2006-01-02 15:04:05"),
			NameList:   make([]string, def.RoomPlayerCount),
			ScoreList:  make([]int32, def.RoomPlayerCount),
			RoundList:  make([]*pf.RecordRound, 0),
		}
		r.RoomRecord = roomRecord
		for _, seat := range r.RoomSeats {
			player := cache.GetPfPlayer(seat.ID)
			if player != nil {
				r.RoomRecord.NameList[seat.Pos] = player.GetNickname()
			}
		}
	}

	r.RoomRecord.RoundList = append(r.RoomRecord.RoundList, &pf.RecordRound{
		CurRound:  r.RoomStatus.CurRound,
		RoundID:   0,
		StartTime: curTime.Format("2006-01-02 15:04:05"),
		ScoreList: make([]int32, def.RoomPlayerCount),
	})

	// DaoRecord
	r.DaoRecord = &dao.Record{
		CreateID:   r.RoomSeats[0].ID,
		CreateTime: r.RoomConfig.CreateTime,
		RoomID:     r.RoomID,
		MjType:     r.RoomConfig.MJType,
		TotalRound: r.RoomConfig.Config.Jushu,
		StartTime:  curTime,
		CurRound:   r.RoomStatus.CurRound,
		EastID:     0,
		SouthID:    0,
		WestID:     0,
		NorthID:    0,
		EastScore:  0,
		SouthScore: 0,
		WestScore:  0,
		NorthScore: 0,
	}
}

func (r *MJRoom) PrintTile(desc string) {
	return
	fmt.Println()
	fmt.Println(desc)
	for _, seat := range r.RoomSeats {
		fmt.Println("Pos", seat.Pos)
		fmt.Println("ID", seat.ID)
		fmt.Println("Score", seat.Score)
		fmt.Println("PaoScore", seat.PaoScore)
		if seat.DrawTile != nil {
			fmt.Println("Draw", string(seat.DrawTile.Unicode))
		}
		if seat.DiscardTile != nil {
			fmt.Println("Draw", string(seat.DiscardTile.Unicode))
		}
		fmt.Println("Hand", seat.HandTiles.GetUnicodeCards())
		fmt.Println("Desk", seat.DeskTiles.GetUnicodeCards())
		fmt.Println("MingGang", seat.MingGangTiles.GetUnicodeCards())
		fmt.Println("AnGang", seat.AnGangTiles.GetUnicodeCards())
		fmt.Println("BuGang", seat.BuGangTiles.GetUnicodeCards())
		fmt.Println("Peng", seat.PengTiles.GetUnicodeCards())
		fmt.Println("Chi", seat.ChiTiles.GetUnicodeCards())
		fmt.Println("IsMingGang", seat.IsMingGangTiles.GetUnicodeCards())
		fmt.Println("IsBuGang", seat.IsBuGangTiles.GetUnicodeCards())
		fmt.Println("IsPeng", seat.IsPengTiles.GetUnicodeCards())
		fmt.Println("IsChi", seat.IsChiTiles.GetUnicodeCards())
		fmt.Println()
	}
	return
}

type MJSeat struct {
	ID         int32 // 玩家id
	Pos        int32 // 玩家位置
	Score      int32 // 玩家分数
	Prepared   bool  // 是否准备
	VoteStatus int32 // 投票状态
	IsBuyPao   bool  // 是否买跑
	PaoScore   int32 // 买跑分数

	ZiMoCnt    int32 // 自摸次数
	JiePaoCnt  int32 // 接炮次数
	DiaoPaoCnt int32 // 点炮次数

	DrawTile    *Tile // 刚抓的牌
	DiscardTile *Tile // 刚打的牌

	HandTiles Tiles // 手里的牌
	DeskTiles Tiles // 桌上的牌

	MingGangTiles Tiles // 明杠的牌
	AnGangTiles   Tiles // 暗杠的牌
	BuGangTiles   Tiles // 补杠的牌
	PengTiles     Tiles // 碰的牌
	ChiTiles      Tiles // 吃的牌

	IsMingGangTiles Tiles // 被明杠的牌
	IsBuGangTiles   Tiles // 被补杠的牌
	IsPengTiles     Tiles // 被碰的牌
	IsChiTiles      Tiles // 被吃的牌

	HuScore   int32 // 胡牌番数
	GangScore int32 // 杠牌番数
}

type MJTiles struct {
	Tiles   Tiles // 整副牌
	DrawPos int32 // 抓牌位置
	LeftCnt int32 // 剩余牌数
	GangCnt int32 // 开杠数目
}

type MJHupai struct {
	IsZiMo       bool    // 是否自摸
	IsZiMoHupai  bool    // 是否自摸胡牌
	ZiMoPos      int32   // 自摸位置
	ZiMoTile     *Tile   // 自摸哪张牌
	ZiMoHuTypes  []int32 // 自摸胡牌类型
	ZiMoHuTiles  Tiles   // 自摸胡牌列表
	ZiMoHuFanshu Fanshu  // 自摸胡牌番数

	IsDianPao       bool    // 是否点炮
	IsDianPaoHupai  bool    // 是否点炮胡牌
	DianPaoPos      int32   // 点炮位置
	JiePaoPos       int32   // 接炮位置
	DianPaoTile     *Tile   // 点炮哪张牌
	DianPaoHuTypes  []int32 // 点炮胡牌类型
	DianPaoHuTiles  Tiles   // 点炮胡牌列表
	DianPaoHuFanshu Fanshu  // 点炮胡牌番数

	IsMultiDianPao       bool      // 是否一炮多响
	IsMultiDianPaoHupai  bool      // 是否一炮多响胡牌
	MultiDianPaoPos      int32     // 点炮位置
	MultiJiePaoPos       []int32   // 接炮位置
	MultiJiePaoHupaiPos  []int32   // 一炮多响胡牌位置
	MultiJiePaoPassPos   []int32   // 一炮多响过位置
	MultiDianPaoTile     *Tile     // 一炮多响哪张牌
	MultiDianPaoHuTypes  [][]int32 // 一炮多响胡牌类型
	MultiDianPaoHuTiles  []Tiles   // 一炮多响胡牌列表
	MultiDianPaoHuFanshu []Fanshu  // 一炮多响胡牌番数
}

type MJStatus struct {
	BankerPos int32 // 庄家位置
	CurPos    int32 // 当前位置
	CurRound  int32 // 当前局数
	LastPos   int32 // 上一个出牌人位置

	buyPaoListRecv      *pf.BuyPaoListRecv      // 断线买跑提示
	operationNoticeRecv *pf.OperationNoticeRecv // 断线重连操作提示
	discardNoticerecv   *pf.DiscardNoticeRecv   // 断线重连出牌提示
	endRoundRecv        *pf.EndRoundRecv        // 断线重连小结算

	GameStatus int32 // 牌局状态
	IsSubCard  bool  // 是否扣除房卡
	IsStart    bool  // 游戏是否开始
	IsGaming   bool  // 是否正在游戏
	IsVoting   bool  // 是否正在投票
	VotePos    int32 // 投票发起人位置
	VoteTime   int32 // 投票剩余时间
}

type MJConfig struct {
	MJType     int32     // 麻将类型
	CreateTime time.Time // 创建时间
	Config     Config    // 牌桌配置
	CliConfig  []int32   // 客户端配置
	Cheat      bool      // 是否开启防作弊 TRUE：开启 FALSE：关闭
}

type MJRoom struct {
	RoomID     int32          // 房间ID
	RoomSeats  []*MJSeat      // 房间玩家
	RoomTiles  *MJTiles       // 房间麻将
	RoomHupai  *MJHupai       // 房间胡牌
	RoomStatus *MJStatus      // 房间状态
	RoomConfig *MJConfig      // 房间配置
	RoomRecord *pf.RecordRoom // 房间战绩
	DaoRecord  *dao.Record    // 每局战绩
}

type RoomMap struct {
	sync.RWMutex
	mmap map[int32]*MJRoom
}

func (m RoomMap) Exist(roomID int32) (ok bool) {
	m.RLock()
	defer m.RUnlock()
	_, ok = m.mmap[roomID]
	return
}

func (m RoomMap) ExistPlayer(playerID, roomID int32) (ok bool) {
	m.RLock()
	defer m.RUnlock()
	room, ok := m.mmap[roomID]
	if ok {
		ok = false
		for _, seat := range room.RoomSeats {
			if seat != nil && seat.ID == playerID {
				ok = true
				break
			}
		}
	}
	return
}

var RoomCache RoomMap

func init() {
	RoomCache = RoomMap{
		mmap: make(map[int32]*MJRoom, def.InitRoomCount),
	}
}
