package enrichers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTmpFile(t *testing.T, path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("error writing file: %s", err)
	}
}

func TestStravaCsvEnrichFailOpen(t *testing.T) {
	_, err := NewStravaCsvEnricher("testdata/missing-file.csv")
	assert.NotNil(t, err)
}

func TestStravaCsvEnrichFailCsvParse(t *testing.T) {
	defer os.RemoveAll("enrich.csv")
	createTmpFile(t, "enrich.csv", []byte("a\n,,\n"))

	_, err := NewStravaCsvEnricher("enrich.csv")
	require.NotNil(t, err)

	assert.Equal(t, "record on line 2: wrong number of fields", err.Error())
}

func TestStravaCsvEnrichActivity(t *testing.T) {
	enricher, err := NewStravaCsvEnricher("testdata/strava_activities.csv")
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
	enricher, err := NewStravaCsvEnricher("testdata/strava_activities.csv")
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
