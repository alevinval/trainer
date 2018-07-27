package enricher

import (
	"testing"

	"github.com/alevinval/trainer/internal/testutil"
	"github.com/alevinval/trainer/internal/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStravaCsvEnrichFailOpen(t *testing.T) {
	_, err := StravaCsv("testdata/missing-file.csv")
	assert.NotNil(t, err)
}

func TestStravaCsvEnrichFailCsvParse(t *testing.T) {
	tmp := testutil.NewTemp()
	filePath := tmp.Create("enrich.csv", []byte("a\n,,\n"))
	defer tmp.Remove()

	_, err := StravaCsv(filePath)
	require.NotNil(t, err)

	assert.Equal(t, "record on line 2: wrong number of fields", err.Error())
}

func TestStravaCsvEnrichActivity(t *testing.T) {
	enricher, err := StravaCsv("testdata/strava_activities.csv")
	require.Nil(t, err)

	a := &trainer.Activity{}
	a.SetMetadata(&trainer.Metadata{
		Name: "",
		DataSource: trainer.DataSource{
			Type: trainer.FileDataSource,
			Name: "783319746.fit.gz",
		},
	})

	enricher.Enrich(a)

	assert.Equal(t, "Lunch Run", a.Metadata().Name)
}

func TestStravaCsvEnrichIgnoresNonFileActivities(t *testing.T) {
	enricher, err := StravaCsv("testdata/strava_activities.csv")
	require.Nil(t, err)

	a := &trainer.Activity{}
	a.SetMetadata(&trainer.Metadata{
		Name: "",
		DataSource: trainer.DataSource{
			Type: trainer.DataSourceType("not-a-file"),
			Name: "783319746.fit.gz",
		},
	})

	enricher.Enrich(a)

	assert.Equal(t, "", a.Metadata().Name)
}
