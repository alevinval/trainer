package trainer

import (
	"fmt"
	"math"
)

type (

	// Cadence of both feet in (steps/min)
	Cadence float64

	// HeartRate in (BPM)
	HeartRate int

	// Pace represents a Speed in (min/km)
	Pace Speed

	// Performance in (steps/s * m/s) / (beats/s)
	Performance float64

	// Speed in (m/s)
	Speed float64

	// Point in terms of geo coordinates
	Point struct {
		Lat, Lon float64
	}
)

func (c Cadence) String() string {
	return fmt.Sprintf("%0.0f steps/s", c)
}

func (h HeartRate) String() string {
	return fmt.Sprintf("%d bpm", h)
}

func (s Pace) String() string {
	totalSecondsFloat := 1000 / float64(s)
	if math.IsInf(totalSecondsFloat, 0) {
		return "n/a"
	}
	totalSeconds := int(totalSecondsFloat)
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%d:%02d min/km", minutes, seconds)
}

func (p Performance) String() string {
	return fmt.Sprintf("%0.2f", p)
}

func (s Speed) String() string {
	return fmt.Sprintf("%0.2f m/s", s)
}

func (p Point) String() string {
	return fmt.Sprintf("lat=%0.6f, lon=%0.6f", p.Lat, p.Lon)
}

func (p Point) distanceTo(other Point) float64 {
	return distance(p.Lat, p.Lon, other.Lat, other.Lon)
}
