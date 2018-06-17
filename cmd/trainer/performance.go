package main

import (
	"log"
	"os"

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
	activities, err := findActivities(path, prefix)
	if err != nil {
		log.Printf("performance command failed: %s", err)
	}

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
