package kf

var pair Tiles
var result Tiles

// InTiles 是否存在
func (t Tiles) InTiles(tile Tile) (ok bool) {
	for _, temp := range t {
		if temp.Unicode == tile.Unicode {
			ok = true
			break
		}
	}
	return
}

// DeletePair 剔除对子
func (t Tiles) DeletePair(i, j int) (tiles Tiles) {
	tiles = make(Tiles, 0, len(t)-2)
	for k := range t {
		if k == i || k == j {
			continue
		}
		tiles = append(tiles, t[k])
	}
	return
}

// SearchPair 查找对子
func (t Tiles) SearchPair() (pairList []Tiles, elseList []Tiles) {
	for i := 0; i < len(t)-1; i++ {
		if t[i].Unicode == t[i+1].Unicode {
			exist := false
			for _, temp := range pairList {
				if temp[0].Unicode == t[i].Unicode {
					exist = true
					break
				}
			}
			if !exist {
				tiles := t.DeletePair(i, i+1)
				pair := Tiles{t[i], t[i+1]}
				pairList = append(pairList, pair)
				elseList = append(elseList, tiles)
			}
		}
	}
	return
}

// IsHupai 胡牌
func (t Tiles) IsHupai(seat *MJSeat, config *MJConfig) (huTiles Tiles, fanshu Fanshu, ok bool) {
	num := len(t)
	tiles := make(Tiles, num)
	copy(tiles, t)
	tiles.SortTiles()

	if len(tiles)%3 != 2 {
		return
	}

	// 七对特殊处理
	if Hu_7pair(tiles) {
		fanshu = Fanshu{
			Qiduifanbei: true,
		}
		huTiles = make(Tiles, len(tiles))
		copy(huTiles, tiles)
		ok = true
		return
	}

	pairList, elseList := tiles.SearchPair()
	for i := range pairList {
		pair = pairList[i]
		result = Tiles{}
		if Hu(elseList[i]) {
			result = append(result, pair...)
			ok = true
			break
		}
	}

	if ok {
		fanshu = Fanshu{}
		huTiles = make(Tiles, len(result))
		copy(huTiles, result)
	}
	return
}

func Hu(tiles Tiles) bool {
	if len(tiles) == 0 {
		return true
	}

	// 判断当前牌有几张
	var count int
	for i := 0; i < len(tiles); i++ {
		if tiles[i].Unicode == tiles[0].Unicode {
			count++
		} else {
			break
		}
	}

	if count >= 4 {
		count = 3
	}
	switch count {
	case 1:
		// 1张想胡牌必须和后面的组成顺子
		newTiles, ok := Hu_1(tiles)
		if !ok {
			return false
		}
		return Hu(newTiles)
	case 2:
		// 2张想胡牌必须和后面的组成顺子
		newTiles, ok := Hu_2(tiles)
		if !ok {
			return false
		}
		return Hu(newTiles)
	case 3:
		// 3张想胡牌可以自己组成刻子也可以和后面组成顺子
		// 3张和后面组成顺子那么也一定有3个刻子
		newTiles, ok := Hu_3(tiles)
		if !ok {
			return false
		}
		return Hu(newTiles)
	case 4:
		// 4张想胡牌可以自己组成刻子也可以和后面组成顺子
		// 4张和后面组成顺子也就变成了7对，而且是至尊超级豪华七对x16
		newTiles, ok := Hu_4(tiles)
		if !ok {
			return false
		}
		return Hu(newTiles)
	}
	return false
}

func Hu_1(tiles Tiles) (left Tiles, ok bool) {
	if tiles[0].Type == Wind {
		return nil, false
	}

	var a1, b1 bool
	var tile_a1, tile_b1 Tile
	for i := 1; i < len(tiles); i++ {
		if tiles[i].Type == tiles[0].Type && tiles[i].Unicode == tiles[0].Unicode+1 && !a1 {
			a1 = true
			tile_a1 = tiles[i]
		} else if tiles[i].Type == tiles[0].Type && tiles[i].Unicode == tiles[0].Unicode+2 && !b1 {
			b1 = true
			tile_b1 = tiles[i]
		} else {
			left = append(left, tiles[i])
		}
	}
	if a1 && b1 {
		result = append(result, tiles[0], tile_a1, tile_b1)
		ok = true
	}
	return
}

func Hu_2(tiles Tiles) (left Tiles, ok bool) {
	if tiles[0].Type == Wind {
		return nil, false
	}

	var a1, b1, a2, b2 bool
	var tile_a1, tile_b1, tile_a2, tile_b2 Tile
	for i := 2; i < len(tiles); i++ {
		if tiles[i].Type == tiles[0].Type && tiles[i].Unicode == tiles[0].Unicode+1 && !a1 {
			a1 = true
			tile_a1 = tiles[i]
		} else if tiles[i].Type == tiles[1].Type && tiles[i].Unicode == tiles[1].Unicode+1 && !a2 {
			a2 = true
			tile_a2 = tiles[i]
		} else if tiles[i].Type == tiles[0].Type && tiles[i].Unicode == tiles[0].Unicode+2 && !b1 {
			b1 = true
			tile_b1 = tiles[i]
		} else if tiles[i].Type == tiles[1].Type && tiles[i].Unicode == tiles[1].Unicode+2 && !b2 {
			b2 = true
			tile_b2 = tiles[i]
		} else {
			left = append(left, tiles[i])
		}
	}
	if a1 && b1 && a2 && b2 {
		result = append(result, tiles[0], tile_a1, tile_b1, tiles[1], tile_a2, tile_b2)
		ok = true
	}
	return
}

func Hu_3(tiles Tiles) (left Tiles, ok bool) {
	for i := 3; i < len(tiles); i++ {
		left = append(left, tiles[i])
	}
	result = append(result, tiles[0], tiles[1], tiles[2])
	ok = true
	return
}

func Hu_4(tiles Tiles) (left Tiles, ok bool) {
	var a1, b1 bool
	var tile_a1, tile_b1 Tile
	for i := 4; i < len(tiles); i++ {
		if tiles[i].Type == tiles[3].Type && tiles[i].Unicode == tiles[3].Unicode+1 && !a1 {
			a1 = true
			tile_a1 = tiles[i]
		} else if tiles[i].Type == tiles[3].Type && tiles[i].Unicode == tiles[3].Unicode+2 && !b1 {
			b1 = true
			tile_b1 = tiles[i]
		} else {
			left = append(left, tiles[i])
		}
	}
	if a1 && b1 {
		result = append(result, tiles[0], tiles[1], tiles[2], tiles[3], tile_a1, tile_b1)
		ok = true
	}
	return
}

func Hu_7pair(tiles Tiles) (ok bool) {
	if len(tiles) == 14 {
		for i := 0; i < 14; i = i + 2 {
			if tiles[i].Unicode != tiles[i+1].Unicode {
				return
			}
		}
		ok = true
	}
	return
}
