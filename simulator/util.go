package simulator

import "math"

type Coordinates struct {
	x  uint
	y  uint
}

func (c Coordinates) Distance(c2 Coordinates) int {
	first := math.Pow(float64(c2.x)-float64(c.x), 2)
	second := math.Pow(float64(c2.y)-float64(c.y), 2)
	return int(math.Sqrt(first + second))
}
