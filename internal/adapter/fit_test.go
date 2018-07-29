package adapter

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tormoder/fit"
)

func TestReadFit(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/sample.fit")
	require.Nil(t, err)

	f, err := Fit(data)
	require.Nil(t, err)

	assert.Equal(t, "", f.Metadata().DataSource.Name)
	assert.Equal(t, trainer.DataSourceType(""), f.Metadata().DataSource.Type)
	assert.Equal(t, "", f.Metadata().Name)
	assert.Equal(t, "2015-01-20 14:04:05 +0000 UTC", f.Metadata().Time.String())
	assert.Equal(t, 260, len(f.DataPoints()))

	// Compare first data point
	expectedTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2015-01-20T13:34:38.000Z")
	expected := &trainer.DataPoint{
		Time: expectedTime,
		Coords: trainer.Point{
			Lat:       25.05563087761402,
			Lon:       121.62827912718058,
			Elevation: 40.200000000000045,
		},
		Hr:    255,
		Cad:   162,
		Speed: 2.69334129650378,
		Perf:  1.7110638824847544,
		N:     1,
	}
	assert.Equal(t, expected, f.DataPoints()[50])
}

func TestReadInvalidFit(t *testing.T) {
	data := []byte("invalid fit data")

	_, err := Fit(data)

	assert.NotNil(t, err)
}

func TestReadNotActivity(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/sample.fit")
	require.Nil(t, err)

	f, err := Fit(data)
	require.Nil(t, err)
	f.(*fitAdapter).file.FileId.Type = fit.FileTypeDevice

	assert.Equal(t, 0, len(f.DataPoints()))
}
