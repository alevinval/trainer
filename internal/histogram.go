package trainer

type (
	Histogram struct {
		data bpmToDataPoints
	}
	FlatHistogram struct {
		data bpmToDataPoint
	}
	bpmToDataPoints map[HeartRate]DataPointList
	bpmToDataPoint  map[HeartRate]*DataPoint
)

func (hist *Histogram) Reset() {
	hist.data = make(bpmToDataPoints)
}

func (hist *Histogram) Feed(provider DataPointProvider) {
	for _, dp := range provider.DataPoints() {
		_, ok := hist.data[dp.Hr]
		if !ok {
			hist.data[dp.Hr] = DataPointList{dp}
		} else {
			hist.data[dp.Hr] = append(hist.data[dp.Hr], dp)
		}
	}
}

func (hist *Histogram) Data() bpmToDataPoints {
	return hist.data
}

func (hist *Histogram) Flatten() *FlatHistogram {
	flat := new(FlatHistogram)
	flat.Reset()
	for bpm, datapoints := range hist.data {
		avg := &DataPoint{
			Speed: datapoints.AvgSpeed(),
			Cad:   datapoints.AvgCad(),
			Perf:  datapoints.AvgPerf(),
			Hr:    bpm,
			N:     len(datapoints),
		}
		flat.data[bpm] = avg
	}
	return flat
}

func (flat *FlatHistogram) Reset() {
	flat.data = make(bpmToDataPoint)
}

func (flat *FlatHistogram) Data() bpmToDataPoint {
	return flat.data
}
