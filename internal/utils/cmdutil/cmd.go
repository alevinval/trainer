package cmdutil

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alevinval/trainer/internal/home"
	"github.com/alevinval/trainer/internal/trainer"
)

// Cluster analysis a set of activities.
func Cluster(activities trainer.ActivityList) {
	for _, cluster := range trainer.GetClusters(activities, trainer.DistanceCriteria(5000.0)) {
		tagCloud := trainer.TagCloudFromActivities(cluster.Activities)
		avgPerf := cluster.Activities.DataPoints().AvgPerf()
		fmt.Printf("%s\n%s\nAvg.perf: %0.2f\n\n", cluster, tagCloud, avgPerf)
	}
}

// Performance analysis for a set of activities.
func Performance(activities trainer.ActivityList, outputPath string) {
	histogram := activities.DataPoints().GetHistogram()
	if len(outputPath) > 0 {
		output, err := os.Create(outputPath)
		if err != nil {
			return
		}
		defer output.Close()
		trainer.WriteCsvTo(histogram, output)
	} else {
		trainer.PrintHistogram(histogram)
		if len(activities) == 1 {
			activity := activities[0]
			fmt.Printf("Overall performance: %s\n", activity.DataPoints().AvgPerf())
			laps := activity.DataPoints().Laps(1000)
			for i, lap := range laps {
				fmt.Printf("Lap %d:\n\t%s (p=%s) (d=%s)\n", i, lap.Pace, lap.Performance, lap.Distance)
			}
		}
	}
}

// Training analysis for a set of activities
func Training(activities trainer.ActivityList, timeWindow string) {
	duration := getTrainingWindowDuration(timeWindow)
	chunks := activities.ChunkByDuration(duration)

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
		if breakTime > duration {
			fmtStr += fmt.Sprintf(" %.1f days break", breakTime.Hours()/24)
		}

		fmtStr += "\n"

		fmt.Printf(fmtStr)

		lastActivity = chunk[len(chunk)-1]
	}
}

func Sync(lookupPath string) {
	home.Sync(lookupPath)
}

func getTrainingWindowDuration(window string) time.Duration {
	switch window {
	case "month":
		return time.Hour * 720
	case "year":
		return time.Hour * 8760
	case "week":
		return time.Hour * 168
	default:
		log.Printf("unrecognized time window: %s", window)
		return time.Hour * 168
	}
}
