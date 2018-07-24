package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	trainingCmd *cobra.Command

	trainingWindow string
	weekWindow     = time.Hour * 168
	monthWindow    = time.Hour * 720
	yearWindow     = time.Hour * 8760
)

func init() {
	trainingCmd = &cobra.Command{
		Use:   "training [path]",
		Short: "display performance evolution across time",
		Long: `Displays evolution of performance across time, shows progress, breaks
		taken, rate of change, etc...`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doTraining(args[0])
		},
	}
	trainingCmd.Flags().StringVar(&trainingWindow, "window", "week", "time-frame used to analyze performance evolution (week, month, year)")
}

func getTrainingWindowDuration() time.Duration {
	switch trainingWindow {
	case "month":
		return monthWindow
	case "year":
		return yearWindow
	}
	return weekWindow
}

func doTraining(path string) {
	activities, err := findActivities(path)
	if err != nil {
		log.Printf("training command failed: %s", err)
	}
	window := getTrainingWindowDuration()
	chunks := activities.ChunkByDuration(window)

	lastActivity := chunks[0][len(chunks[0])-1]
	for _, chunk := range chunks {
		dateFromStr := chunk[0].Metadata().Time.Format("2006-01-02")
		fmtStr := fmt.Sprintf("%s: p=%s hr=%d cad=%.0f speed=%0.2f [a=%-3d n=%-5d]",
			dateFromStr,
			chunk.DataPoints().AvgPerf(),
			chunk.DataPoints().AvgHeartRate(),
			chunk.DataPoints().AvgCad(),
			chunk.DataPoints().AvgSpeed(),
			len(chunk),
			len(chunk.DataPoints()),
		)

		breakTime := chunk[0].Metadata().Time.Sub(lastActivity.Metadata().Time)
		if breakTime > window {
			fmtStr += fmt.Sprintf(" %.1f days break", breakTime.Hours()/24)
		}

		fmtStr += "\n"

		fmt.Printf(fmtStr)

		lastActivity = chunk[len(chunk)-1]
	}
}
