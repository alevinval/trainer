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

func findPathsWithPrefix(root string, prefix string) (filePaths chan string, err error) {
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

func findActivities(lookupPath, prefix string) (activities trainer.ActivityList, err error) {
	activities = trainer.ActivityList{}

	paths, err := findPathsWithPrefix(lookupPath, prefix)
	if err != nil {
		log.Printf("cannot find activities in %s: %s\n", lookupPath, err)
		return nil, err
	}

	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
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
	err = applyEnrichers(activities)
	if err != nil {
		return
	}
	if filterByName != "" {
		activities = activities.Filter(func(a *trainer.Activity) bool {
			cloud := trainer.TagCloudFromActivities(trainer.ActivityList{a})
			return cloud.Contains(filterByName)
		})
	}
	return
}

func applyEnrichers(activities trainer.ActivityList) (err error) {
	list := []trainer.Enricher{}
	if enrichStravaCsvPath != "" {
		e, err := enrichers.NewStravaCsvEnricher(enrichStravaCsvPath)
		if err != nil {
			log.Printf("cannot apply  enricher: %s", err)
			return err
		}
		list = append(list, e)
	}
	return trainer.EnrichActivities(activities, list...)
}
