package trainer

type (
	Histogram struct {
		Name string
		data bpmToDataPoints
	}
	bpmToDataPoints map[HeartRate]DataPointList
)

func (hist *Histogram) Reset() {
	hist.data = make(bpmToDataPoints)
}

func (hist *Histogram) Feed(provider dataPointProvider) {
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

func (hist *Histogram) Flatten() *Histogram {
	flat := new(Histogram)
	flat.Reset()
	for bpm, datapoints := range hist.data {
		avg := &DataPoint{
			Speed: datapoints.AvgSpeed(),
			Cad:   datapoints.AvgCad(),
			Perf:  datapoints.AvgPerf(),
			Hr:    bpm,
			n:     len(datapoints),
		}
		flat.Feed(avg)
	}
	return flat
}

func (hist *Histogram) GetAvgPerf() Performance {
	list := DataPointList{}
	for _, datapoints := range hist.data {
		list = append(list, datapoints...)
	}
	return list.AvgPerf()
}

func (hist *Histogram) maxCountHeartRate() (max int) {
	for _, list := range hist.data {
		if len(list) > max {
			max = len(list)
		}
	}
	return
}
