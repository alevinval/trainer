package main

import (
	"log"
	"os"
	"sync"

	"github.com/alevinval/trainer"
	"github.com/spf13/cobra"
)

var (
	performanceCmd    *cobra.Command
	performanceOutput string
)

func init() {
	performanceCmd = &cobra.Command{
		Use:   "performance [path]",
		Short: "compute performance data for the matched activities",
		Long: `computes performance data and builds a histogram to analyse how
		you perform on each heart rate zone.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doPerformanceCommand(args[0])
		},
	}
	performanceCmd.Flags().StringVar(&performanceOutput, "output", "", "file name to output processed data")
}

func doPerformanceCommand(path string) error {
	activities := findActivities(path, prefix)
	histogram := activities.DataPoints().GetHistogram()
	if len(performanceOutput) > 0 {
		output, err := os.Create(performanceOutput)
		if err != nil {
			return err
		}
		defer output.Close()
		trainer.WriteCsvTo(histogram, output)
	} else {
		trainer.PrintHistogram(histogram)
	}
	return nil
}

func findActivities(lookupPath, prefix string) trainer.ActivityList {
	activities := trainer.ActivityList{}

	fileNames, err := findFilesWithPrefix(lookupPath, prefix)
	if err != nil {
		log.Printf("cannot find activities in %s: %s\n", lookupPath, err)
		return activities
	}

	wg := new(sync.WaitGroup)
	mux := new(sync.Mutex)
	for fileName := range fileNames {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			activity, err := trainer.OpenFile(fileName)
			if err != nil {
				log.Printf("cannot open file %q: %s\n", fileName, err)
				return
			}
			mux.Lock()
			activities = append(activities, activity)
			mux.Unlock()
		}(fileName)
	}
	wg.Wait()
	return activities
}
