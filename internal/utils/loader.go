package utils

import (
	"log"
	"sync"

	"github.com/alevinval/trainer/internal/provider"
	"github.com/alevinval/trainer/internal/trainer"
)

// ActivitiesFromPaths allocates a fixed size pool of workers to load
// activities from a list of paths.
func ActivitiesFromPaths(paths []string) (list trainer.ActivityList) {
	wg := new(sync.WaitGroup)
	wg.Add(len(paths))

	inputCh := make(chan string, len(paths))
	for i := range paths {
		inputCh <- paths[i]
	}
	close(inputCh)

	activitiesCh := make(chan trainer.ActivityProvider, len(paths))

	maxParallelOpen := 10
	for w := 0; w < maxParallelOpen; w++ {
		go activityLoader(wg, inputCh, activitiesCh)
	}
	wg.Wait()
	close(activitiesCh)

	list = make(trainer.ActivityList, 0)
	for activity := range activitiesCh {
		list = append(list, activity)
	}
	return
}

func activityLoader(wg *sync.WaitGroup, paths <-chan string, activities chan<- trainer.ActivityProvider) {
	for path := range paths {
		defer wg.Done()
		provider, err := provider.File(path)
		if err != nil {
			log.Printf("cannot open file %q: %s\n", path, err)
			continue
		}
		activities <- provider
	}
}
