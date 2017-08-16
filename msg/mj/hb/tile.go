package hb

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
	Front   int32 // 前端展示使用
	IsHun   bool  // 是否混牌
	IsHu    bool  // 是否胡牌
}

type Tiles []Tile

type TilesList []Tiles

// Unicode enum
const (
	// Prevailing wind tiles
	EastWind int32 = iota + 0x1F000
	SouthWind
	WestWind
	NorthWind

	// Dragon tiles
	RedDragon   // hongzhong
	GreenDragon // qingfa
	WhiteDragon // baiban

	// Character suit tiles
	OneOfCharacters // wan
	TwoOfCharacters
	ThreeOfCharacters
	FourOfCharacters
	FiveOfCharacters
	SixOfCharacters
	SevenOfCharacters
	EightOfCharacters
	NineOfCharacters

	// Bamboo suit tiles
	OneOfBamboos // tiao
	TwoOfBamboos
	ThreeOfBamboos
	FourOfBamboos
	FiveOfBamboos
	SixOfBamboos
	SevenOfBamboos
	EightOfBamboos
	NineOfBamboos

	// Circle suit tiles
	OneOfCircles // bing
	TwoOfCircles
	ThreeOfCircles
	FourOfCircles
	FiveOfCircles
	SixOfCircles
	SevenOfCircles
	EightOfCircles
	NineOfCircles

	// Flower tiles
	PlumFlower          // mei
	OrchidFlower        // lan
	BambooFlower        // zhu
	ChrysanthemumFlower // ju

	// Season tiles
	SpringSeason
	SummerSeason
	AutumnSeason
	WinterSeason

	// Miscellaneous tiles
	JokerMiscellaneous // baida
	BackMiscellaneous
)
