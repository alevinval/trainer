package main

import (
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"

	"github.com/alevinval/trainer"
	"github.com/alevinval/trainer/enrichers"
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
	return
}

func getPathsWithPrefix(root string, prefix string) (filePaths chan string, err error) {
	paths, err := findFiles(root)
	if err != nil {
		return
	}
	prefix = path.Join(root, prefix)
	filePaths = make(chan string)
	go func() {
		for filePath := range paths {
			if len(filePath) < len(prefix) {
				continue
			}
			if strings.Compare(filePath[:len(prefix)], prefix) != 0 {
				continue
			}
			filePaths <- filePath
		}
		close(filePaths)
	}()
	return
}

func getActivitiesFromPaths(paths chan string) trainer.ActivityList {
	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)

	activities := trainer.ActivityList{}
	for filePath := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			activity, err := trainer.OpenFile(path)
			if err != nil {
				log.Printf("cannot open file %q: %s\n", path, err)
				return
			}
			mux.Lock()
			activities = append(activities, activity)
			mux.Unlock()
		}(filePath)
	}
	wg.Wait()
	return activities
}

func findFiles(root string) (filePaths chan string, err error) {
	paths, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	filePaths = make(chan string)
	go func() {
		for _, filePath := range paths {
			filePaths <- path.Join(root, filePath.Name())
		}
		close(filePaths)
	}()
	return
}

func applyStravaEnricher(activities trainer.ActivityList) (err error) {
	stravaEnricher, err := enrichers.NewStravaCsvEnricher(stravaCsvEnrichPath)
	if err != nil {
		log.Printf("cannot apply  enricher: %s", err)
		return err
	}
	return trainer.EnrichActivities(activities, stravaEnricher)
}
