package trainer

import (
	"fmt"
	"io"
	"math"
	"sort"
)

func (hist *Histogram) WriteTo(w io.Writer) {
	sortedBpm := []int{}
	for bpm := range hist.Data() {
		sortedBpm = append(sortedBpm, int(bpm))
	}
	sort.Ints(sortedBpm)

	flat := hist.Flatten()
	fmt.Fprint(w, "BPM,N,Pace,Speed,Cadence,Perf(steps/s * m/s / bps)\n")
	for _, bpm := range sortedBpm {
		datapoints := hist.Data()[HeartRate(bpm)]
		dp := flat.Data()[HeartRate(bpm)][0]
		fmt.Fprintf(w, "%d,%d,%0.2f,%0.2f,%0.0f,%0.2f\n", bpm, len(datapoints), Pace(dp.Speed), dp.Speed, dp.Cad, dp.Perf)
	}
}

func getMaxWidth(hist *Histogram) (max int) {
	for _, datapoints := range hist.Data() {
		if len(datapoints) > max {
			max = len(datapoints)
		}
	}
	return
}

func (hist *Histogram) PrintRaw() {
	hrArr := []int{}
	for hr := range hist.Data() {
		hrArr = append(hrArr, int(hr))
	}
	sort.Ints(hrArr)

	maxWidth := float64(getMaxWidth(hist))

	flat := hist.Flatten()
	for _, hrInt := range hrArr {
		hr := HeartRate(hrInt)
		datapoints := hist.Data()[hr]
		dots := ""
		width := math.Floor(float64(len(datapoints)) / maxWidth * 50)
		for width > 0 {
			width--
			dots += "·"
		}
		if len(dots) == 0 {
			continue
		}
		bar := fmt.Sprintf("%s | p=%s | %s\n", hr, flat.Data()[hr][0].Perf, dots)
		fmt.Print(bar)
	}
}
