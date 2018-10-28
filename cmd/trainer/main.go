package main

import (
	"log"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/alevinval/trainer/internal/utils/cmdutil"
	"github.com/spf13/cobra"
)

var cmdArgs cmdutil.CmdArgs

func init() {
	cmdArgs = cmdutil.CmdArgs{}
}

func loadActivities(lookupPath string) trainer.ActivityList {
	cmdArgs.LookupPath = lookupPath
	activities, err := cmdutil.LoadActivityFromArgs(cmdArgs)
	if err != nil {
		log.Fatalf("error loading activities: %s", err)
	}
	return activities
}

func main() {
	log.SetFlags(0)

	var (
		cmdTrainingWindow    string
		cmdPerformanceOutput string
	)

	cmd := cobra.Command{Use: "trainer"}

	clusterCmd := &cobra.Command{
		Use:   "cluster [path]",
		Short: "clusters activities by coordinates and computes their performance",
		Long: `clusters activities by coordinates and computes the performance histograms
			for each of them.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activities := loadActivities(args[0])
			cmdutil.Cluster(activities)
		},
	}

	performanceCmd := &cobra.Command{
		Use:   "performance [path]",
		Short: "compute performance data for the matched activities",
		Long: `computes performance data and builds a histogram to analyse how
		you perform on each heart rate zone.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activities := loadActivities(args[0])
			cmdutil.Performance(activities, cmdPerformanceOutput)
		},
	}

	trainingCmd := &cobra.Command{
		Use:   "training [path]",
		Short: "display performance evolution across time",
		Long: `Displays evolution of performance across time, shows progress, breaks
		taken, rate of change, etc...`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			activities := loadActivities(args[0])
			cmdutil.Training(activities, cmdTrainingWindow)
		},
	}

	cmd.PersistentFlags().StringVar(&cmdArgs.FilterByPrefix, "prefix", "", "only process files matching the prefix")
	cmd.PersistentFlags().StringVar(&cmdArgs.FilterByName, "name", "", "Filters activities whose name does not match with the filter")
	cmd.PersistentFlags().StringVar(&cmdArgs.FilterByDate, "date", "", "Filters activities whose date does not match with the filter")
	cmd.PersistentFlags().StringVar(&cmdArgs.FilterByDateFrom, "date-from", "", "Filters activities whose date is above the specified date prefix (inclusive)")
	cmd.PersistentFlags().StringVar(&cmdArgs.FilterByDateTo, "date-to", "", "Filters activities whose date is below the specified date prefix (non-inclusive)")
	cmd.PersistentFlags().BoolVar(&cmdArgs.LogDebug, "debug", false, "Log debug traces")
	cmd.PersistentFlags().StringVar(&cmdArgs.StravaCsvEnrichPath, "strava-csv-enrich", "", "enrich metadata from a csv file")

	performanceCmd.Flags().StringVar(&cmdPerformanceOutput, "output", "", "file name to output processed data")

	trainingCmd.Flags().StringVar(&cmdTrainingWindow, "window", "week", "time-frame used to analyze performance evolution (week, month, year)")

	cmd.AddCommand(performanceCmd)
	cmd.AddCommand(clusterCmd)
	cmd.AddCommand(trainingCmd)
	cmd.Execute()
}
