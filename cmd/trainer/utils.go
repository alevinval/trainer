package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/alevinval/trainer/internal"
	"github.com/alevinval/trainer/internal/enrichers"
	"github.com/alevinval/trainer/pkg/providers"
)

func findActivities(lookupPath string) (activities trainer.ActivityList, err error) {
	paths, err := getPathsWithPrefix(lookupPath, filterByPrefix)
	if err != nil {
		log.Printf("cannot find activities in %s: %s\n", lookupPath, err)
		return nil, err
	}

	activities = getActivitiesFromPaths(paths)

	if stravaCsvEnrichPath != "" {
		err = applyStravaEnricher(activities)
		if err != nil {
			log.Printf("cannot apply strava csv enricher: %s", err)
			return
		}
	}

	if filterByDate != "" {
		activities = activities.Filter(func(a *trainer.Activity) bool {
			date := a.Metadata().Time.Format("20060102")
			return strings.HasPrefix(date, filterByDate)
		})
	}

	if filterByDateFrom != "" {
		activities = activities.Filter(func(a *trainer.Activity) bool {
			date := a.Metadata().Time.Format("20060102")
			return date[0:len(filterByDateFrom)] >= filterByDateFrom
		})
	}

	if filterByDateTo != "" {
		activities = activities.Filter(func(a *trainer.Activity) bool {
			date := a.Metadata().Time.Format("20060102")
			return date[0:len(filterByDateTo)] < filterByDateTo
		})
	}

	if filterByName != "" {
		activities = activities.Filter(func(a *trainer.Activity) bool {
			cloud := trainer.TagCloudFromActivities(trainer.ActivityList{a})
			return cloud.Contains(filterByName)
		})
	}

	activities.SortByTime()

	if logDebug {
		for _, activity := range activities {
			log.Printf("Activity: %s", activity.Metadata().Name)
		}
	}

	return
}

func getPathsWithPrefix(root string, prefix string) (prefixedPaths []string, err error) {
	prefix = path.Join(root, prefix)

	paths, err := findPaths(root)
	if err != nil {
		return nil, err
	}

	prefixedPaths = make([]string, 0)
	for i := range paths {
		if len(paths[i]) < len(prefix) {
			continue
		}
		if strings.Compare(paths[i][:len(prefix)], prefix) != 0 {
			continue
		}
		prefixedPaths = append(prefixedPaths, paths[i])
	}
	return
}

func getActivitiesFromPaths(paths []string) (list trainer.ActivityList) {
	wg := new(sync.WaitGroup)
	wg.Add(len(paths))

	inputCh := make(chan string, len(paths))
	for i := range paths {
		inputCh <- paths[i]
	}
	close(inputCh)

	activitiesCh := make(chan *trainer.Activity, len(paths))

	maxParallelOpen := 10
	for w := 0; w < maxParallelOpen; w++ {
		go loadActivityWorker(wg, inputCh, activitiesCh)
	}
	wg.Wait()
	close(activitiesCh)

	list = make(trainer.ActivityList, 0)
	for activity := range activitiesCh {
		list = append(list, activity)
	}
	return
}

func loadActivityWorker(wg *sync.WaitGroup, paths <-chan string, activities chan<- *trainer.Activity) {
	for path := range paths {
		defer wg.Done()
		activity, err := providers.OpenFile(path)
		if err != nil {
			log.Printf("cannot open file %q: %s\n", path, err)
			continue
		}
		activities <- activity
	}
}

func findPaths(root string) (filePaths []string, err error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	filePaths = make([]string, len(files))
	for i := range files {
		filePaths[i] = path.Join(root, files[i].Name())
	}
	return
}

func applyStravaEnricher(activities trainer.ActivityList) (err error) {
	stravaEnricher, err := enrichers.NewStravaCsvEnricher(stravaCsvEnrichPath)
	if err != nil {
		return err
	}
	return trainer.EnrichActivities(activities, stravaEnricher)
}
