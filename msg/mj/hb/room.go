package hb

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
	if config.Config.Daihuner {
		hunCnt := tiles.CountHunCard(config.HunTile.Front)
		if hunCnt > 0 {
			huTiles, fanshu, ok = tiles.IsHupaiWithHun(s, config, int(hunCnt))
		} else {
			huTiles, fanshu, ok = tiles.IsHupai(s, config)
			if ok && config.Config.Suhu {
				fanshu.Suhu = true
			}
		}
	} else {
		huTiles, fanshu, ok = tiles.IsHupai(s, config)
	}
	return
}

// IsDianPao 点炮
func (s *MJSeat) IsDianPao(config *MJConfig, dianpaoTile *Tile) (huTiles Tiles, fanshu Fanshu, ok bool) {
	tiles := make(Tiles, len(s.HandTiles))
	copy(tiles, s.HandTiles)
	dianpaoTile.IsHu = true
	if config.Config.Daihuner {
		hunCnt := tiles.CountHunCard(config.HunTile.Front)
		tiles = append(tiles, *dianpaoTile)
		if hunCnt > 0 {
			huTiles, fanshu, ok = tiles.IsHupaiWithHun(s, config, int(hunCnt))
		} else {
			huTiles, fanshu, ok = tiles.IsHupai(s, config)
			if ok && config.Config.Suhu {
				fanshu.Suhu = true
			}
		}
	} else {
		tiles = append(tiles, *dianpaoTile)
		huTiles, fanshu, ok = tiles.IsHupai(s, config)
	}
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

// IsChi 吃
func (s *MJSeat) IsChi(tile Tile) (chiTypes []int32) {
	tiles := s.HandTiles
	if len(tiles) < 2 {
		return
	}

	if tile.Type == Wind {
		return
	}

	pTileTwo := Tile{
		Type:    tile.Type,
		Unicode: tile.Unicode - 2,
	}
	pTileOne := Tile{
		Type:    tile.Type,
		Unicode: tile.Unicode - 1,
	}
	nTileOne := Tile{
		Type:    tile.Type,
		Unicode: tile.Unicode + 1,
	}
	nTileTwo := Tile{
		Type:    tile.Type,
		Unicode: tile.Unicode + 2,
	}

	var ispTileTwo bool
	var ispTileOne bool
	var isnTileOne bool
	var isnTileTwo bool

	for i := 0; i < len(tiles); i++ {
		if tiles[i].Type == pTileTwo.Type && tiles[i].Unicode == pTileTwo.Unicode {
			ispTileTwo = true
		}
		if tiles[i].Type == pTileOne.Type && tiles[i].Unicode == pTileOne.Unicode {
			ispTileOne = true
		}
		if tiles[i].Type == nTileOne.Type && tiles[i].Unicode == nTileOne.Unicode {
			isnTileOne = true
		}
		if tiles[i].Type == nTileTwo.Type && tiles[i].Unicode == nTileTwo.Unicode {
			isnTileTwo = true
		}
	}

	if ispTileTwo && ispTileOne {
		chiTypes = append(chiTypes, def.ChiType1)
	}
	if ispTileOne && isnTileOne {
		chiTypes = append(chiTypes, def.ChiType2)
	}
	if isnTileOne && isnTileTwo {
		chiTypes = append(chiTypes, def.ChiType3)
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
		fanshu.Zimo = true
		fanshu.Menqing = s.IsMenqing()
		fanshu.Gangshangkaihua = afterGang
		fanshu.Haidilaoyue = room.RoomTiles.IsLiuju()
		fanshu.Daizhuangxian = room.RoomConfig.Config.Daizhuangxian
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
			fanshu.Menqing = s.IsMenqing()
			fanshu.Daizhuangxian = room.RoomConfig.Config.Daizhuangxian
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

func (s *MJSeat) GetOtherChiList(tile Tile, pos int32) (operationList []*pf.OperationItem) {
	// 吃
	if s.Pos == (pos+1)%def.RoomPlayerCount {
		chiTypes := s.IsChi(tile)
		for _, chiType := range chiTypes {
			cardsList := []int32{}
			switch chiType {
			case def.ChiType1:
				cardsList = []int32{tile.Front - 2, tile.Front - 1, tile.Front}
			case def.ChiType2:
				cardsList = []int32{tile.Front - 1, tile.Front, tile.Front + 1}
			case def.ChiType3:
				cardsList = []int32{tile.Front, tile.Front + 1, tile.Front + 2}
			}
			operationList = append(operationList, &pf.OperationItem{
				Type:      def.Chi,
				Keycard:   tile.Front,
				CardsList: cardsList,
				SrcPos:    pos + 1,
				DesPos:    s.Pos + 1,
			})
		}
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

func (s *MJSeat) AddGangTypes(huTypes *[]int32) {
	for i := 0; i < len(s.IsMingGangTiles); i++ {
		*huTypes = append(*huTypes, 22)
	}
	for i := 0; i < len(s.AnGangTiles); i += 4 {
		*huTypes = append(*huTypes, 23)
	}
	for i := 0; i < len(s.MingGangTiles); i += 4 {
		*huTypes = append(*huTypes, 24)
	}
	for i := 0; i < len(s.BuGangTiles); i += 4 {
		*huTypes = append(*huTypes, 25)
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
		fanshu, huTypes := r.RoomHupai.ZiMoHuFanshu.GetScore()
		r.RoomHupai.ZiMoHuTypes = huTypes
		if r.RoomConfig.Config.Daizhuangxian {
			if r.RoomHupai.ZiMoPos == r.RoomStatus.BankerPos {
				fanshu = 2 * fanshu
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.ZiMoPos {
						seat.HuScore = (def.RoomPlayerCount - 1) * fanshu
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			} else {
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.ZiMoPos {
						seat.HuScore = def.RoomPlayerCount * fanshu
					} else if seat.Pos == r.RoomStatus.BankerPos {
						seat.HuScore = -2 * fanshu
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			}
		} else {
			for _, seat := range r.RoomSeats {
				if seat.Pos == r.RoomHupai.ZiMoPos {
					seat.HuScore = (def.RoomPlayerCount - 1) * fanshu
				} else {
					seat.HuScore = -1 * fanshu
				}
			}
		}
	} else if r.RoomHupai.IsDianPaoHupai {
		fanshu, huTypes := r.RoomHupai.DianPaoHuFanshu.GetScore()
		r.RoomHupai.DianPaoHuTypes = huTypes
		if r.RoomConfig.Config.Dianpaoyijiachu {
			if r.RoomConfig.Config.Daizhuangxian {
				if r.RoomStatus.BankerPos == r.RoomHupai.DianPaoPos || r.RoomStatus.BankerPos == r.RoomHupai.JiePaoPos {
					fanshu = 2 * fanshu
				}
			}
			r.RoomSeats[r.RoomHupai.DianPaoPos].HuScore = -1 * fanshu
			r.RoomSeats[r.RoomHupai.JiePaoPos].HuScore = fanshu
		} else if r.RoomConfig.Config.Dianpaosanjiachu {
			if r.RoomHupai.JiePaoPos == r.RoomStatus.BankerPos {
				if r.RoomConfig.Config.Daizhuangxian {
					fanshu = 2 * fanshu
				}
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						seat.HuScore = def.RoomPlayerCount * fanshu
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						seat.HuScore = -2 * fanshu
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			} else if r.RoomHupai.DianPaoPos == r.RoomStatus.BankerPos {
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = (def.RoomPlayerCount + 2) * fanshu
						} else {
							seat.HuScore = def.RoomPlayerCount * fanshu
						}
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = -4 * fanshu
						} else {
							seat.HuScore = -2 * fanshu
						}
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			} else {
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = (def.RoomPlayerCount + 1) * fanshu
						} else {
							seat.HuScore = def.RoomPlayerCount * fanshu
						}
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						seat.HuScore = -2 * fanshu
					} else if seat.Pos == r.RoomStatus.BankerPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = -2 * fanshu
						} else {
							seat.HuScore = -1 * fanshu
						}
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			}
		} else if r.RoomConfig.Config.Dianpaodabao {
			if r.RoomConfig.Config.Daizhuangxian {
				if r.RoomStatus.BankerPos == r.RoomHupai.DianPaoPos || r.RoomStatus.BankerPos == r.RoomHupai.JiePaoPos {
					fanshu = 2 * fanshu
				}
			}
			r.RoomSeats[r.RoomHupai.DianPaoPos].HuScore = -3 * fanshu
			r.RoomSeats[r.RoomHupai.JiePaoPos].HuScore = 3 * fanshu
		}
	} else if r.RoomHupai.IsMultiDianPaoHupai {
		// for _, jiePaoPos := range r.RoomHupai.MultiJiePaoHupaiPos {
		// 	for i, pos := range r.RoomHupai.MultiJiePaoPos {
		// 		if pos == jiePaoPos {
		// 			break
		// 		}
		// 	}
		// }

		fanshu, huTypes := r.RoomHupai.DianPaoHuFanshu.GetScore()
		r.RoomHupai.DianPaoHuTypes = huTypes
		if r.RoomConfig.Config.Dianpaoyijiachu {
			if r.RoomConfig.Config.Daizhuangxian {
				if r.RoomStatus.BankerPos == r.RoomHupai.DianPaoPos || r.RoomStatus.BankerPos == r.RoomHupai.JiePaoPos {
					fanshu = 2 * fanshu
				}
			}
			r.RoomSeats[r.RoomHupai.DianPaoPos].HuScore = -1 * fanshu
			r.RoomSeats[r.RoomHupai.JiePaoPos].HuScore = fanshu
		} else if r.RoomConfig.Config.Dianpaosanjiachu {
			if r.RoomHupai.JiePaoPos == r.RoomStatus.BankerPos {
				if r.RoomConfig.Config.Daizhuangxian {
					fanshu = 2 * fanshu
				}
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						seat.HuScore = def.RoomPlayerCount * fanshu
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						seat.HuScore = -2 * fanshu
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			} else if r.RoomHupai.DianPaoPos == r.RoomStatus.BankerPos {
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = (def.RoomPlayerCount + 2) * fanshu
						} else {
							seat.HuScore = def.RoomPlayerCount * fanshu
						}
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = -4 * fanshu
						} else {
							seat.HuScore = -2 * fanshu
						}
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			} else {
				for _, seat := range r.RoomSeats {
					if seat.Pos == r.RoomHupai.JiePaoPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = (def.RoomPlayerCount + 1) * fanshu
						} else {
							seat.HuScore = def.RoomPlayerCount * fanshu
						}
					} else if seat.Pos == r.RoomHupai.DianPaoPos {
						seat.HuScore = -2 * fanshu
					} else if seat.Pos == r.RoomStatus.BankerPos {
						if r.RoomConfig.Config.Daizhuangxian {
							seat.HuScore = -2 * fanshu
						} else {
							seat.HuScore = -1 * fanshu
						}
					} else {
						seat.HuScore = -1 * fanshu
					}
				}
			}
		} else if r.RoomConfig.Config.Dianpaodabao {
			if r.RoomConfig.Config.Daizhuangxian {
				if r.RoomStatus.BankerPos == r.RoomHupai.DianPaoPos || r.RoomStatus.BankerPos == r.RoomHupai.JiePaoPos {
					fanshu = 2 * fanshu
				}
			}
			r.RoomSeats[r.RoomHupai.DianPaoPos].HuScore = -3 * fanshu
			r.RoomSeats[r.RoomHupai.JiePaoPos].HuScore = 3 * fanshu
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
				huTypes = r.RoomHupai.ZiMoHuTypes
				seat.ZiMoCnt++
				drawCard = r.RoomHupai.ZiMoTile.Front
				seat.HandTiles = r.RoomHupai.ZiMoHuTiles.DelHuCard()
			}
			seat.AddGangTypes(&huTypes)
		} else if r.RoomHupai.IsDianPaoHupai {
			if seat.Pos == r.RoomHupai.JiePaoPos {
				isHu = 1
				huTypes = r.RoomHupai.DianPaoHuTypes
				seat.JiePaoCnt++
				drawCard = r.RoomHupai.DianPaoTile.Front
				seat.HandTiles = r.RoomHupai.ZiMoHuTiles.DelHuCard()
			} else if seat.Pos == r.RoomHupai.DianPaoPos {
				seat.DiaoPaoCnt++
				huTypes = append(huTypes, 21)
			}
			seat.AddGangTypes(&huTypes)
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
				huTypes = r.RoomHupai.MultiDianPaoHuTypes[index]
				seat.JiePaoCnt++
				drawCard = r.RoomHupai.MultiDianPaoTile.Front
				seat.HandTiles = r.RoomHupai.MultiDianPaoHuTiles[index].DelHuCard()
			} else if seat.Pos == r.RoomHupai.MultiDianPaoPos {
				seat.DiaoPaoCnt++
				huTypes = append(huTypes, 21)
			}
			seat.AddGangTypes(&huTypes)
		} else {
			if r.RoomConfig.Config.Huangzhuanghuanggang {
				seat.GangScore = 0
			}
			if !r.RoomConfig.Config.Huangzhuanghuanggang {
				seat.AddGangTypes(&huTypes)
			}
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
			Coins:      -1 * coins,
			Cards:      -1 * cards,
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
	for _, seat := range r.RoomSeats {
		if seat.Pos == desPos {
			seat.GangScore += 2 * (def.RoomPlayerCount - 1)
		} else {
			seat.GangScore -= 2
		}
	}
}

func (r *MJRoom) AddMingGangScore(srcPos int32, desPos int32) {
	for _, seat := range r.RoomSeats {
		if seat.Pos == desPos {
			seat.GangScore += 2
		} else if seat.Pos == srcPos {
			seat.GangScore -= 2
		}
	}
}

func (r *MJRoom) AddBuGangScore(desPos int32) {
	for _, seat := range r.RoomSeats {
		if seat.Pos == desPos {
			seat.GangScore += 1 * (def.RoomPlayerCount - 1)
		} else {
			seat.GangScore -= 1
		}
	}
}

func (r *MJRoom) ClearSeats() {
	for _, seat := range r.RoomSeats {
		seat.Prepared = false

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
	fmt.Println()
	fmt.Println(desc)
	for _, seat := range r.RoomSeats {
		fmt.Println("Pos", seat.Pos)
		fmt.Println("ID", seat.ID)
		fmt.Println("Score", seat.Score)
		if r.RoomConfig.HunTile != nil {
			fmt.Println("HunTile", string(r.RoomConfig.HunTile.Unicode))
		}
		if seat.DrawTile != nil {
			fmt.Println("Draw", string(seat.DrawTile.Unicode))
		}
		if seat.DiscardTile != nil {
			fmt.Println("Discard", string(seat.DiscardTile.Unicode))
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
	HunTile    *Tile     // 混牌
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
