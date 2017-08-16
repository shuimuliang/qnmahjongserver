package slice

// +gen slice:"Aggregate[string]"
type Employee struct {
	Name       string
	Department string
}

// +gen slice:"All"
type Person struct {
	Name    string
	Present bool
}

// +gen slice:"Any"
type Person2 struct {
	Name    string
	Present bool
}

// +gen slice:"Average"
type Celsius float64

// +gen slice:"Average[int]"
type Player struct {
	Name   string
	Points int
}

// +gen slice:"Count"
type Monster struct {
	Name  string
	Furry bool
	Fangs int
}

// +gen slice:"Distinct"
type Hipster struct {
	FavoriteBand string
	Mustachioed  bool
	Bepectacled  bool
}

// +gen slice:"DistinctBy"
type Hipster2 struct {
	FavoriteBand string
	Mustachioed  bool
}

// +gen slice:"First"
type Customer struct {
	Name string
	Here bool
}

// +gen slice:"GroupBy[int]"
type Movie struct {
	Title string
	Year  int
}

// +gen slice:"Max"
type Price float64

// +gen slice:"Max[Dollars]"
type Movie2 struct {
	Title     string
	BoxOffice Dollars
}

type Dollars int

// +gen slice:"MaxBy"
type Rectangle struct {
	Width, Height int
}

// +gen slice:"Min"
type Price2 float64

// +gen slice:"Min[Dollars]"
type Movie3 struct {
	Title     string
	BoxOffice Dollars
}

// +gen slice:"MinBy"
type Rectangle2 struct {
	Width, Height int
}

// +gen slice:"Select[int]"
type Player2 struct {
	Name   string
	Points int
}

// +gen slice:"Shuffle"
type Rating int

// +gen slice:"Sort,SortDesc"
type Rating2 int

// +gen slice:"SortBy"
type Movie4 struct {
	Title string
	Year  int
}

// +gen slice:"Where"
type Movie5 struct {
	Title string
	Year  int
}

// +gen slice:"Shuffle,SortBy"
type Tile struct {
	Type    int32
	Value   int32
	Unicode int32
	Front   int32 // 前端展示使用
}
