package adapter

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGpx(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/sample.gpx")
	require.Nil(t, err)

	g, err := Gpx(data)
	require.Nil(t, err)

	assert.Equal(t, "", g.Metadata().DataSource.Name)
	assert.Equal(t, trainer.DataSourceType(""), g.Metadata().DataSource.Type)
	assert.Equal(t, "Some activity name", g.Metadata().Name)
	assert.Equal(t, "2015-01-20 13:26:30 +0000 UTC", g.Metadata().Time.String())
	assert.Equal(t, 3, len(g.DataPoints()))

	// Compare first data point
	expectedTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2017-06-19T16:49:40.000Z")
	expected := &trainer.DataPoint{
		Time: expectedTime,
		Coords: trainer.Point{
			Lat:       1.000001,
			Lon:       1.000001,
			Elevation: 100,
		},
		Hr:    94,
		Cad:   170,
		Speed: 0.03148350908743588,
		Perf:  0.05693826111557553,
		N:     1,
	}
	assert.Equal(t, expected, g.DataPoints()[1])
}

func TestReadInvalidGpx(t *testing.T) {
	data := []byte("invalid gpx data")

	_, err := Gpx(data)
	assert.NotNil(t, err)
}
