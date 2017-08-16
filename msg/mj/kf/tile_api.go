package kf

import "sort"

func (tiles Tiles) GetFrontCards() (cards []int32) {
	cards = make([]int32, len(tiles))
	for i, tile := range tiles {
		cards[i] = tile.Front
	}
	return
}

func (tiles Tiles) GetUnicodeCards() (cards []string) {
	temps := make(Tiles, len(tiles))
	copy(temps, tiles)
	temps.SortCards()
	cards = make([]string, len(tiles))
	for i, tile := range temps {
		cards[i] = string(tile.Unicode)
	}
	return
}

func (tiles Tiles) GetEndTile(front int32) (tile Tile) {
	for _, tile = range tiles {
		if tile.Front == front {
			break
		}
	}
	return
}

func (tiles Tiles) GetEndTiles(fronts []int32) (result Tiles) {
	result = make(Tiles, 0, len(fronts))
	for _, front := range fronts {
		result = append(result, tiles.GetEndTile(front))
	}
	return
}

func (tiles Tiles) DelFrontCard(front int32) (result Tiles) {
	del := false
	result = make(Tiles, 0, len(tiles))
	for _, tile := range tiles {
		if !del && tile.Front == front {
			del = true
			continue
		}
		result = append(result, tile)
	}
	return
}

func (tiles Tiles) DelHuCard() (result Tiles) {
	for _, tile := range tiles {
		if tile.IsHu {
			continue
		}
		result = append(result, tile)
	}
	return
}

func (tiles Tiles) DelPengCard(front int32) (result Tiles) {
	result = tiles.DelFrontCard(front)
	result = result.DelFrontCard(front)
	return
}

func (tiles Tiles) DelGangCard(front int32) (result Tiles) {
	result = make(Tiles, 0, len(tiles))
	for _, tile := range tiles {
		if tile.Front == front {
			continue
		}
		result = append(result, tile)
	}
	return
}

func (tiles Tiles) DelChiCard(chiTiles Tiles) (result Tiles) {
	result = tiles.DelFrontCard(chiTiles[0].Front)
	result = result.DelFrontCard(chiTiles[1].Front)
	return
}

func (tiles Tiles) SortTiles() {
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Unicode <= tiles[j].Unicode
	})
}

func (tiles Tiles) SortCards() {
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Front <= tiles[j].Front
	})
}
