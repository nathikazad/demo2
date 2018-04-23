package simulator

import "math"

type Coordinates struct {
	X  uint
	Y  uint
}

func (c Coordinates) Distance(c2 Coordinates) float64 {
	first := math.Pow(float64(c2.X)-float64(c.X), 2)
	second := math.Pow(float64(c2.Y)-float64(c.Y), 2)
	return math.Sqrt(first + second)
}
