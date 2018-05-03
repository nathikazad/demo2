package sim2

import (
  "math"
)

// coordmath - Descrbes float-based coordinates with utility transforms

// Coords - Structure for generic 2D coordinates.
type Coords struct {
  X float64
  Y float64
}

// Distance - Calculate the distance between two coordinates.
func (c1 Coords) Distance(c2 Coords) float64 {
	sqdx := math.Pow(c2.X - c1.X, 2)
	sqdy := math.Pow(c2.Y - c1.Y, 2)
	return math.Sqrt(sqdx + sqdy)
}

// Angle - Calculate the angle in degrees formed between coordinate vector and x-axis.
func (c1 Coords) Angle(c2 Coords) float64 {
	dx := c2.X - c1.X
	dy := c2.Y - c1.Y
  return math.Atan2(dy, dx) * 180 / math.Pi
}

// UnitVector - Calculate the unit vector (starting at c1) created from c1 to c2.
func (c1 Coords) UnitVector(c2 Coords) Coords {
  // Find vector components from c1 -> c2
  dx := c2.X - c1.X
  dy := c2.Y - c1.Y

  // Find magnitude of vector from c1 -> c2
  mag := c1.Distance(c2)

  // Find unit vector components from c1 -> c2
  uvecx := dx / mag
  uvecy := dy / mag

  return Coords{ X:uvecx, Y:uvecy }
}

// ProjectInDirection - project a new pair of coordinates some distance from current position.
//   Positive numUnits indicates projecting towards c2 from c2; negative is alternate (away).
//   Note that the ending coordinates may be negative in value if crossing quadrants.
func (c1 Coords) ProjectInDirection(numUnits float64, c2 Coords) (targ Coords) {
  // Find unit vector from c1 -> c2
  uvec := c1.UnitVector(c2)

  // Scale unit vector by desired amplitude
  targ.X = uvec.X * numUnits
  targ.Y = uvec.Y * numUnits
  return
}
