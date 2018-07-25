package enrichers

import (
	"encoding/csv"
	"os"
	"path"

	"github.com/alevinval/trainer/internal"
)

type stravaCsvEnricher struct {
	activityNameCol        int
	fileNameCol            int
	fileNameToActivityName map[string]string
}

// NewStravaCsvEnricher returns an enricher of activites
// using strava dumps metadata.
func NewStravaCsvEnricher(filePath string) (trainer.Enricher, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	enricher := &stravaCsvEnricher{
		activityNameCol:        2,
		fileNameCol:            9,
		fileNameToActivityName: map[string]string{},
	}
	enricher.buildActivityNameMap(records)
	return enricher, nil
}

func (e *stravaCsvEnricher) buildActivityNameMap(records [][]string) {
	for _, columns := range records {
		fileName := path.Base(columns[e.fileNameCol])
		activityName := columns[e.activityNameCol]
		e.fileNameToActivityName[fileName] = activityName
	}
}

func (e *stravaCsvEnricher) Enrich(a *trainer.Activity) (err error) {
	m := a.Metadata()

	// Only enrich activities coming from files.
	if m.DataSource.Type != trainer.FileDataSource {
		return
	}

	// See if the strava metadata csv contains activity name
	// for that file
	fileName := path.Base(m.DataSource.Name)
	activityName, hasName := e.fileNameToActivityName[fileName]
	if hasName {
		a.Metadata().Name = activityName
	}

	return
}
