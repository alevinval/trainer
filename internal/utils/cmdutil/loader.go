package cmdutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/alevinval/trainer/internal/enricher"
	"github.com/alevinval/trainer/internal/trainer"
	"github.com/alevinval/trainer/internal/utils"
)

// CmdArgs wraps together common arguments that commands accept to
// filter, enrich and select which kinds of activities must be loaded.
type CmdArgs struct {
	LookupPath          string
	StravaCsvEnrichPath string
	FilterByPrefix      string
	FilterByName        string
	FilterByDate        string
	FilterByDateFrom    string
	FilterByDateTo      string
	LogDebug            bool
	HomePath            string
}

// LoadActivityFromArgs loads activities applying filters and enrichers as needed.
func LoadActivityFromArgs(args CmdArgs) (activities trainer.ActivityList, err error) {
	paths, err := getPathsWithPrefix(args.LookupPath, args.FilterByPrefix)
	if err != nil {
		log.Printf("cannot find activities in %s: %s\n", args.LookupPath, err)
		return nil, err
	}

	activities = utils.ActivitiesFromPaths(paths)

	if args.StravaCsvEnrichPath != "" {
		err = applyStravaEnricher(activities, args.StravaCsvEnrichPath)
		if err != nil {
			log.Printf("cannot apply strava csv enricher: %s", err)
			return
		}
	}

	if args.FilterByDate != "" {
		activities = activities.Filter(func(a trainer.ActivityProvider) bool {
			date := a.Metadata().Time.Format("20060102")
			return strings.HasPrefix(date, args.FilterByDate)
		})
	}

	if args.FilterByDateFrom != "" {
		activities = activities.Filter(func(a trainer.ActivityProvider) bool {
			date := a.Metadata().Time.Format("20060102")
			return date[0:len(args.FilterByDateFrom)] >= args.FilterByDateFrom
		})
	}

	if args.FilterByDateTo != "" {
		activities = activities.Filter(func(a trainer.ActivityProvider) bool {
			date := a.Metadata().Time.Format("20060102")
			return date[0:len(args.FilterByDateTo)] < args.FilterByDateTo
		})
	}

	if args.FilterByName != "" {
		activities = activities.Filter(func(a trainer.ActivityProvider) bool {
			cloud := trainer.TagCloudFromActivities(trainer.ActivityList{a})
			return cloud.Contains(args.FilterByName)
		})
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("no activities found")
	}

	activities.SortByTime()

	if args.LogDebug {
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

func applyStravaEnricher(activities trainer.ActivityList, stravaCsvEnrichPath string) (err error) {
	stravaEnricher, err := enricher.StravaCsv(stravaCsvEnrichPath)
	if err != nil {
		return err
	}
	return trainer.EnrichActivities(activities, stravaEnricher)
}
