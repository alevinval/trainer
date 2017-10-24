package trainer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	speed3m30s      = Speed(4.76)
	speedStationary = Speed(0)
)

func TestUnitsToString(t *testing.T) {
	for _, test := range []struct {
		unit     fmt.Stringer
		expected string
	}{
		{unit: Cadence(179), expected: "179 steps/s"},
		{unit: Pace(speed3m30s), expected: "3:30 min/km"},
		{unit: Pace(0), expected: "n/a"},
		{unit: Performance(1.12345), expected: "1.12"},
		{unit: Point{1, 2}, expected: "lat=1.000000, lon=2.000000"},
		{unit: HeartRate(129), expected: "129 bpm"},
		{unit: speed3m30s, expected: "4.76 m/s"},
	} {
		assert.Equal(t, test.unit.String(), test.expected)
	}
}
