package trainer

import "math"

const earthRadius = 6378100
const radiansFactor = math.Pi / 180

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1 *= radiansFactor
	lon1 *= radiansFactor
	lat2 *= radiansFactor
	lon2 *= radiansFactor
	h := hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*hsin(lon2-lon1)
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}
