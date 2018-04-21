package main

import "math"

func (p Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.x)-float64(p.x), 2)
	second := math.Pow(float64(p2.y)-float64(p.y), 2)
	return math.Sqrt(first + second)
}
