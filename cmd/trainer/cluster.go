package main

import (
	"fmt"

	"github.com/alevinval/trainer"
	"github.com/spf13/cobra"
)

var clusterCmd *cobra.Command

func init() {
	clusterCmd = &cobra.Command{
		Use:   "cluster [path]",
		Short: "clusters activities by coordinates and computes their performance",
		Long: `clusters activities by coordinates and computes the performance histograms
		for each of them.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doClusterCommand(args[0])
		},
	}
}

func doClusterCommand(path string) error {
	activities := findActivities(path, prefix)
	for _, cluster := range activities.GetClusters() {
		tagCloud := trainer.TagCloudFromActivities(cluster.Activities)
		avgPerf := cluster.Activities.DataPoints().AvgPerf()
		fmt.Printf("%s\n%s\nAvg.perf: %0.2f\n\n", cluster, tagCloud, avgPerf)
	}
	return nil
}
