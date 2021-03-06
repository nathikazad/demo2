package simulator

import (
	"math"
)

type Coordinates struct {
	X  uint
	Y  uint
}

func (c1 Coordinates) Distance(c2 Coordinates) float64 {
	first := math.Pow(float64(c2.X)-float64(c1.X), 2)
	second := math.Pow(float64(c2.Y)-float64(c1.Y), 2)
	return math.Sqrt(first + second)
}

func (c1 Coordinates) Angle(c2 Coordinates) float64 {
	first := float64(c2.X)-float64(c1.X)
	second := float64(c2.Y)-float64(c1.Y)
	return math.Atan2(second, first) * 180 / math.Pi
}

func determineOrientation (currentAngle float64, desiredAngle float64, steps float64) float64 {
	error := (desiredAngle-currentAngle)
	newAngle := 0.0
	if math.Abs(error)  > steps {
		if error > 0 {
			if error < 180 {
				newAngle = currentAngle + steps
			} else {
				newAngle = currentAngle - steps
			}
		} else {
			if error > -180 {
				newAngle = currentAngle - steps
			} else {
				newAngle = currentAngle + steps
			}
		}
	} else {
		newAngle = desiredAngle
	}
	if newAngle > 360 {
		newAngle -= 360
	} else if newAngle < 0 {
		newAngle += 360
	}
	return newAngle

}

// ProjectInDirection project a new pair of coordinates some distance from current position
// positive numUnits indicates projecting towards c2 from c1, negative indicates projecting away
// If the results are outside of the positive quadrant, return zero
// TODO: consider returning closest intersect on the positive axes
func (c1 Coordinates) ProjectInDirection(numUnits int, c2 Coordinates) Coordinates {
  var target Coordinates

  // Find vector components from c1 -> c2
  dX := float64(c2.X) - float64(c1.X)
  dY := float64(c2.Y) - float64(c1.Y)

  // Find magnitude of vector from c1 -> c2
  mag := c1.Distance(c2)

  // Find unit vector components from c1 -> c2
  uVecX := dX / mag
  uVecY := dY / mag

  // Scale unit vector components by desired magnitude
  // Note that this preserves the sign of the magnitude (and therefore the direction)
  targdX := uVecX * float64(numUnits)
  targdY := uVecY * float64(numUnits)

  // Calculate endpoints from starting coordinates
  endX := math.Round(float64(c1.X) + targdX)
  endY := math.Round(float64(c1.Y) + targdY)

  // Bound within positive quadrant
  if endX < 0 || endY < 0 {
    target.X = 0
    target.Y = 0
  } else {
    target.X = uint(endX)
    target.Y = uint(endY)
  }
  //fmt.Print("dX: ", dX, ", dY: ", dY, ", mag: ", mag, ", uVecX: ", uVecX)
  //fmt.Print(", uVecY: ", uVecY, ", targdX: ", targdX, ", targdY: ", targdY)
  //fmt.Print(", endX: ", endX, ", endY: ", endY)
  //fmt.Println(", target.X: ", target.X, ", target.Y: ", target.Y)
  return target
}
