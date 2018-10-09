package enricher

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/alevinval/trainer/internal/trainer"
)

type stravaCsvEnricher struct {
	activityNameCol        int
	fileNameCol            int
	fileNameToActivityName map[string]string
}

// StravaCsv returns an enricher of activites
// using strava dumps metadata.
func StravaCsv(filePath string) (trainer.Enricher, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
		fileName := filepath.Base(columns[e.fileNameCol])
		activityName := columns[e.activityNameCol]
		e.fileNameToActivityName[fileName] = activityName
	}
}

func (e *stravaCsvEnricher) Enrich(provider trainer.ActivityProvider) (err error) {
	m := provider.Metadata()

	// Only enrich activities coming from files.
	if m.DataSource.Type != trainer.FileDataSource {
		return
	}

	// See if the strava metadata csv contains activity name
	// for that file
	fileName := filepath.Base(m.DataSource.Name)
	activityName, hasName := e.fileNameToActivityName[fileName]
	if hasName {
		m.Name = activityName
	}

	return
}
