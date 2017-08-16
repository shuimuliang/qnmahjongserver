package util

import (
	"math"
)

const (
	EARTH_RADIUS = 6371000
)

type Point struct {
	Lat float64
	Lng float64
}

func (p *Point) GreatCircleDistance(p2 *Point) int32 {
	dLat := (p2.Lat - p.Lat) * (math.Pi / 180.0)
	dLon := (p2.Lng - p.Lng) * (math.Pi / 180.0)

	Lat1 := p.Lat * (math.Pi / 180.0)
	Lat2 := p2.Lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(Lat1) * math.Cos(Lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int32(EARTH_RADIUS * c)
}
