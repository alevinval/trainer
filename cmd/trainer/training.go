package main

import (
	"fmt"
	"log"
	"strings"
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

func getTrainingWindowDuration() (time.Duration, string) {
	switch trainingWindow {
	case "month":
		return monthWindow, "month"
	case "year":
		return yearWindow, "year"
	}
	return weekWindow, "week"
}

func doTraining(path string) {
	activities, err := findActivities(path)
	if err != nil {
		log.Printf("training command failed: %s", err)
	}
	window, windowStr := getTrainingWindowDuration()
	chunks := activities.ChunkByDuration(window)

	lastActivity := chunks[0][len(chunks[0])-1]
	for i, chunk := range chunks {
		fmtStr := fmt.Sprintf("%s %2d: p=%s a=%-3d n=%-5d",
			strings.Title(windowStr),
			i,
			chunk.DataPoints().AvgPerf(),
			len(chunk),
			len(chunk.DataPoints()),
		)

		breakTime := chunk[0].Metadata().Time.Sub(lastActivity.Metadata().Time)
		if breakTime > window {
			fmtStr += fmt.Sprintf(" [break=%.1f days]", breakTime.Hours()/24)
		}

		fmtStr += "\n"

		fmt.Printf(fmtStr)

		lastActivity = chunk[len(chunk)-1]
	}
}
