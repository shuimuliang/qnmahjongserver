package tile

import "sync"

const (
	Wind       = 1 // 风牌
	Characters = 2 // 万牌
	Bamboos    = 3 // 条牌
	Circles    = 4 // 筒牌
)

const (
	East  = 1
	South = 2
	West  = 3
	North = 4
	Red   = 5
	Green = 6
	White = 7
)

const (
	One   = 1
	Two   = 2
	Three = 3
	Four  = 4
	Five  = 5
	Six   = 6
	Seven = 7
	Eight = 8
	Nine  = 9
)

type Tile struct {
	Type    int32
	Value   int32
	Unicode int32
	Front   int32
	IsHun   bool
	IsHu    bool
}

type Tiles []Tile

type TilesList []Tiles

var tilesWithFeng = Tiles{
	{Wind, East, 0x1F000, 31, false, false},
	{Wind, South, 0x1F001, 33, false, false},
	{Wind, West, 0x1F002, 32, false, false},
	{Wind, North, 0x1F003, 34, false, false},
	{Wind, Red, 0x1F004, 10, false, false},
	{Wind, Green, 0x1F005, 20, false, false},
	{Wind, White, 0x1F006, 30, false, false},

	{Wind, East, 0x1F000, 31, false, false},
	{Wind, South, 0x1F001, 33, false, false},
	{Wind, West, 0x1F002, 32, false, false},
	{Wind, North, 0x1F003, 34, false, false},
	{Wind, Red, 0x1F004, 10, false, false},
	{Wind, Green, 0x1F005, 20, false, false},
	{Wind, White, 0x1F006, 30, false, false},

	{Wind, East, 0x1F000, 31, false, false},
	{Wind, South, 0x1F001, 33, false, false},
	{Wind, West, 0x1F002, 32, false, false},
	{Wind, North, 0x1F003, 34, false, false},
	{Wind, Red, 0x1F004, 10, false, false},
	{Wind, Green, 0x1F005, 20, false, false},
	{Wind, White, 0x1F006, 30, false, false},

	{Wind, East, 0x1F000, 31, false, false},
	{Wind, South, 0x1F001, 33, false, false},
	{Wind, West, 0x1F002, 32, false, false},
	{Wind, North, 0x1F003, 34, false, false},
	{Wind, Red, 0x1F004, 10, false, false},
	{Wind, Green, 0x1F005, 20, false, false},
	{Wind, White, 0x1F006, 30, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},
}

var tilesWithoutFeng = Tiles{
	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Characters, One, 0x1F007, 1, false, false},
	{Characters, Two, 0x1F008, 2, false, false},
	{Characters, Three, 0x1F009, 3, false, false},
	{Characters, Four, 0x1F00A, 4, false, false},
	{Characters, Five, 0x1F00B, 5, false, false},
	{Characters, Six, 0x1F00C, 6, false, false},
	{Characters, Seven, 0x1F00D, 7, false, false},
	{Characters, Eight, 0x1F00E, 8, false, false},
	{Characters, Nine, 0x1F00F, 9, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Circles, One, 0x1F010, 11, false, false},
	{Circles, Two, 0x1F011, 12, false, false},
	{Circles, Three, 0x1F012, 13, false, false},
	{Circles, Four, 0x1F013, 14, false, false},
	{Circles, Five, 0x1F014, 15, false, false},
	{Circles, Six, 0x1F015, 16, false, false},
	{Circles, Seven, 0x1F016, 17, false, false},
	{Circles, Eight, 0x1F017, 18, false, false},
	{Circles, Nine, 0x1F018, 19, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},

	{Bamboos, One, 0x1F019, 21, false, false},
	{Bamboos, Two, 0x1F01A, 22, false, false},
	{Bamboos, Three, 0x1F01B, 23, false, false},
	{Bamboos, Four, 0x1F01C, 24, false, false},
	{Bamboos, Five, 0x1F01D, 25, false, false},
	{Bamboos, Six, 0x1F01E, 26, false, false},
	{Bamboos, Seven, 0x1F01F, 27, false, false},
	{Bamboos, Eight, 0x1F020, 28, false, false},
	{Bamboos, Nine, 0x1F021, 29, false, false},
}

var tilesHunMap = map[int32]Tile{
	30: {Wind, East, 0x1F000, 31, false, false},
	31: {Wind, South, 0x1F001, 33, false, false},
	33: {Wind, West, 0x1F002, 32, false, false},
	32: {Wind, North, 0x1F003, 34, false, false},
	34: {Wind, Red, 0x1F004, 10, false, false},
	10: {Wind, Green, 0x1F005, 20, false, false},
	20: {Wind, White, 0x1F006, 30, false, false},

	9: {Characters, One, 0x1F007, 1, false, false},
	1: {Characters, Two, 0x1F008, 2, false, false},
	2: {Characters, Three, 0x1F009, 3, false, false},
	3: {Characters, Four, 0x1F00A, 4, false, false},
	4: {Characters, Five, 0x1F00B, 5, false, false},
	5: {Characters, Six, 0x1F00C, 6, false, false},
	6: {Characters, Seven, 0x1F00D, 7, false, false},
	7: {Characters, Eight, 0x1F00E, 8, false, false},
	8: {Characters, Nine, 0x1F00F, 9, false, false},

	19: {Circles, One, 0x1F010, 11, false, false},
	11: {Circles, Two, 0x1F011, 12, false, false},
	12: {Circles, Three, 0x1F012, 13, false, false},
	13: {Circles, Four, 0x1F013, 14, false, false},
	14: {Circles, Five, 0x1F014, 15, false, false},
	15: {Circles, Six, 0x1F015, 16, false, false},
	16: {Circles, Seven, 0x1F016, 17, false, false},
	17: {Circles, Eight, 0x1F017, 18, false, false},
	18: {Circles, Nine, 0x1F018, 19, false, false},

	29: {Bamboos, One, 0x1F019, 21, false, false},
	21: {Bamboos, Two, 0x1F01A, 22, false, false},
	22: {Bamboos, Three, 0x1F01B, 23, false, false},
	23: {Bamboos, Four, 0x1F01C, 24, false, false},
	24: {Bamboos, Five, 0x1F01D, 25, false, false},
	25: {Bamboos, Six, 0x1F01E, 26, false, false},
	26: {Bamboos, Seven, 0x1F01F, 27, false, false},
	27: {Bamboos, Eight, 0x1F020, 28, false, false},
	28: {Bamboos, Nine, 0x1F021, 29, false, false},
}

var tilesFrontMap = map[int32]Tile{
	31: {Wind, East, 0x1F000, 31, false, false},
	33: {Wind, South, 0x1F001, 33, false, false},
	32: {Wind, West, 0x1F002, 32, false, false},
	34: {Wind, North, 0x1F003, 34, false, false},
	10: {Wind, Red, 0x1F004, 10, false, false},
	20: {Wind, Green, 0x1F005, 20, false, false},
	30: {Wind, White, 0x1F006, 30, false, false},

	1: {Characters, One, 0x1F007, 1, false, false},
	2: {Characters, Two, 0x1F008, 2, false, false},
	3: {Characters, Three, 0x1F009, 3, false, false},
	4: {Characters, Four, 0x1F00A, 4, false, false},
	5: {Characters, Five, 0x1F00B, 5, false, false},
	6: {Characters, Six, 0x1F00C, 6, false, false},
	7: {Characters, Seven, 0x1F00D, 7, false, false},
	8: {Characters, Eight, 0x1F00E, 8, false, false},
	9: {Characters, Nine, 0x1F00F, 9, false, false},

	11: {Circles, One, 0x1F010, 11, false, false},
	12: {Circles, Two, 0x1F011, 12, false, false},
	13: {Circles, Three, 0x1F012, 13, false, false},
	14: {Circles, Four, 0x1F013, 14, false, false},
	15: {Circles, Five, 0x1F014, 15, false, false},
	16: {Circles, Six, 0x1F015, 16, false, false},
	17: {Circles, Seven, 0x1F016, 17, false, false},
	18: {Circles, Eight, 0x1F017, 18, false, false},
	19: {Circles, Nine, 0x1F018, 19, false, false},

	21: {Bamboos, One, 0x1F019, 21, false, false},
	22: {Bamboos, Two, 0x1F01A, 22, false, false},
	23: {Bamboos, Three, 0x1F01B, 23, false, false},
	24: {Bamboos, Four, 0x1F01C, 24, false, false},
	25: {Bamboos, Five, 0x1F01D, 25, false, false},
	26: {Bamboos, Six, 0x1F01E, 26, false, false},
	27: {Bamboos, Seven, 0x1F01F, 27, false, false},
	28: {Bamboos, Eight, 0x1F020, 28, false, false},
	29: {Bamboos, Nine, 0x1F021, 29, false, false},
}

type TilesSlice struct {
	sync.RWMutex
	tilesWithFeng    Tiles
	tilesWithoutFeng Tiles
}

type TilesMap struct {
	sync.RWMutex
	tilesHunMap   map[int32]Tile
	tilesFrontMap map[int32]Tile
}

var tilesSlice TilesSlice
var tilesMap TilesMap
