package enrichers

import (
	"encoding/csv"
	"os"
	"path"

	"github.com/alevinval/trainer"
)

type stravaCsvEnricher struct {
	r               *csv.Reader
	fileNameCol     int
	activityNameCol int

	fileToActivityMap map[string]string
}

func NewStravaCsvEnricher(filePath string) (enricher *stravaCsvEnricher, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	enricher = &stravaCsvEnricher{
		r:                 csv.NewReader(file),
		fileNameCol:       9,
		activityNameCol:   2,
		fileToActivityMap: map[string]string{},
	}
	records, err := enricher.r.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, columns := range records {
		fileName := path.Base(columns[enricher.fileNameCol])
		activityName := columns[enricher.activityNameCol]
		enricher.fileToActivityMap[fileName] = activityName
	}
	return
}

func (e *stravaCsvEnricher) Enrich(a *trainer.Activity) (err error) {
	if a.Metadata().DataSource.Type != trainer.FileDataSource {
		return nil
	}
	fileName := path.Base(a.Metadata().DataSource.Name)
	activityName, ok := e.fileToActivityMap[fileName]
	if !ok {
		return
	}
	enrichedMetadata := a.Metadata()
	enrichedMetadata.Name = activityName
	a.SetMetadata(enrichedMetadata)
	return nil
}
