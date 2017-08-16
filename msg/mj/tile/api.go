package tiles

import "sort"

func init() {
	tilesSlice = TilesSlice{
		tilesWithFeng:    tilesWithFeng,
		tilesWithoutFeng: tilesWithoutFeng,
	}
	tilesMap = TilesMap{
		tilesHunMap:   tilesHunMap,
		tilesFrontMap: tilesFrontMap,
	}
}

func GetRoomTiles(withFeng bool) Tiles {
	tilesSlice.RLocker()
	defer tilesSlice.RUnlock()

	if withFeng {
		tiles := make(Tiles, len(tilesSlice.tilesWithFeng))
		copy(tiles, tilesSlice.tilesWithFeng)
		return tiles
	}

	tiles := make(Tiles, len(tilesSlice.tilesWithoutFeng))
	copy(tiles, tilesSlice.tilesWithoutFeng)
	return tiles
}

func GetHunTile(front int32) Tile {
	tilesMap.RLock()
	defer tilesMap.RUnlock()

	return tilesMap.tilesHunMap[front]
}

func GetEndTile(front int32) Tile {
	tilesMap.RLock()
	defer tilesMap.RUnlock()

	return tilesMap.tilesFrontMap[front]
}

func GetEndTiles(frontList []int32) Tiles {
	tilesMap.RLock()
	defer tilesMap.RUnlock()

	tiles := make(Tiles, len(frontList))
	for i, front := range frontList {
		tiles[i] = tilesMap.tilesFrontMap[front]
	}
	return tiles
}

func (tiles Tiles) GetFrontTiles() []int32 {
	list := make([]int32, len(tiles))
	for i, tile := range tiles {
		list[i] = tile.Front
	}
	return list
}

func (tiles Tiles) GetUnicodeTiles() []string {
	list := make([]string, len(tiles))
	for i, t := range tiles {
		list[i] = string(t.Unicode)
	}
	return list
}

func (tiles Tiles) DelHuTile() Tiles {
	newTiles := make(Tiles, len(tiles)-1)
	i := 0
	for _, tile := range tiles {
		if !tile.IsHu {
			newTiles[i] = tile
			i++
		}
	}
	return newTiles
}

func (tiles Tiles) DelHunTiles() Tiles {
	count := tiles.CountHunTile()
	newTiles := make(Tiles, len(tiles)-int(count))
	i := 0
	for _, tile := range tiles {
		if !tile.IsHun {
			newTiles[i] = tile
			i++
		}
	}
	return newTiles
}

func (tiles Tiles) CountHunTile() int32 {
	count := int32(0)
	for _, tile := range tiles {
		if tile.IsHun {
			count++
		}
	}
	return count
}

func (tiles Tiles) SortTiles() {
	sort.Slice(tiles, func(i, j int) bool {
		return tiles[i].Unicode <= tiles[j].Unicode
	})
}
