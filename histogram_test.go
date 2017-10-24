package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bpm130 = HeartRate(130)
	bpm150 = HeartRate(150)
	bpm180 = HeartRate(180)
)

var input = DataPointList{
	&DataPoint{Hr: bpm130},
	&DataPoint{Hr: bpm150},
	&DataPoint{Hr: bpm150},
	&DataPoint{Hr: bpm180},
	&DataPoint{Hr: bpm180},
	&DataPoint{Hr: bpm180},
}

func TestHistogramFeed(t *testing.T) {
	h := new(Histogram)
	h.Reset()
	h.Feed(input)

	for _, test := range []struct {
		hr    HeartRate
		count int
	}{
		{bpm130, 1},
		{bpm150, 2},
		{bpm180, 3},
	} {
		list := h.Data()[test.hr]
		assert.Equal(t, test.count, len(list))
	}
}
