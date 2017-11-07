package trainer

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

func getMaxWidth(hist *Histogram) (max int) {
	for _, datapoints := range hist.Data() {
		if datapoints[0].n > max {
			max = datapoints[0].n
		}
	}
	return
}

func getSortedHeartRate(hist *Histogram) []HeartRate {
	list := make([]HeartRate, 0, len(hist.Data()))
	for hr := range hist.Data() {
		list = append(list, hr)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return list
}

func WriteCsvTo(hist *Histogram, w io.Writer) {
	flat := hist.Flatten()
	hrs := getSortedHeartRate(hist)
	fmt.Fprint(w, "BPM,N,Pace,Speed,Cadence,Perf(steps/s * m/s / bps)\n")
	for _, hr := range hrs {
		dp := flat.Data()[hr][0]
		fmt.Fprintf(w, "%d,%d,%0.2f,%0.2f,%0.0f,%0.2f\n", hr, dp.n, Pace(dp.Speed), dp.Speed, dp.Cad, dp.Perf)
	}
}

func PrintHistogram(hist *Histogram) {
	flat := hist.Flatten()
	maxWidth := getMaxWidth(flat)
	hrs := getSortedHeartRate(hist)
	maxDots := 50
	for _, hr := range hrs {
		datapoint := flat.Data()[hr][0]
		numDots := int(math.Floor(50 * float64(datapoint.n) / float64(maxWidth)))
		if numDots == 0 {
			continue
		}
		dots := strings.Repeat(".", numDots)
		dots += strings.Repeat(" ", maxDots-numDots)
		fmt.Printf("%3d bpm | p=%0.2f | %s |\n", hr, datapoint.Perf, dots)
	}
}
