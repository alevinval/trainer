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

func WriteCsvTo(hist *Histogram, w io.Writer) {
	flat := hist.Flatten()

	bpms := []int{}
	for bpm := range flat.Data() {
		bpms = append(bpms, int(bpm))
	}
	sort.Ints(bpms)

	fmt.Fprint(w, "BPM,N,Pace,Speed,Cadence,Perf(steps/s * m/s / bps)\n")
	for _, bpm := range bpms {
		dp := flat.Data()[HeartRate(bpm)][0]
		fmt.Fprintf(w, "%d,%d,%0.2f,%0.2f,%0.0f,%0.2f\n", bpm, dp.n, Pace(dp.Speed), dp.Speed, dp.Cad, dp.Perf)
	}
}

func PrintHistogram(hist *Histogram) {
	flat := hist.Flatten()
	maxWidth := getMaxWidth(hist)

	bpms := []int{}
	for bpm := range flat.Data() {
		bpms = append(bpms, int(bpm))
	}
	sort.Ints(bpms)

	for _, bpm := range bpms {
		hr := HeartRate(bpm)
		datapoints := flat.Data()[hr]
		width := int(math.Floor(float64(datapoints[0].n) / float64(maxWidth)))
		if width == 0 {
			continue
		}
		dots := strings.Repeat(".", width)
		bar := fmt.Sprintf("%s | p=%s | %s\n", hr, flat.Data()[hr][0].Perf, dots)
		fmt.Print(bar)
	}
}
