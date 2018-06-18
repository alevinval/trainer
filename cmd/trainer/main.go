package main

import (
	"github.com/spf13/cobra"
)

var (
	stravaCsvEnrichPath string
	filterByPrefix      string
	filterByName        string
	filterByDate        string
	filterByDateFrom    string
	filterByDateTo      string
	cmd                 *cobra.Command
)

func init() {
	cmd = &cobra.Command{Use: "trainer"}
	cmd.PersistentFlags().StringVar(&stravaCsvEnrichPath, "strava-csv-enrich", "", "enrich metadata from a csv file")
	cmd.PersistentFlags().StringVar(&filterByPrefix, "prefix", "", "only process files matching the prefix")
	cmd.PersistentFlags().StringVar(&filterByName, "name", "", "Filters activities whose name does not match with the filter")
	cmd.PersistentFlags().StringVar(&filterByDate, "date", "", "Filters activities whose date does not match with the filter")
	cmd.PersistentFlags().StringVar(&filterByDateFrom, "date-from", "", "Filters activities whose date is above the specified date prefix (inclusive)")
	cmd.PersistentFlags().StringVar(&filterByDateTo, "date-to", "", "Filters activities whose date is below the specified date prefix (non-inclusive)")
}

func main() {
	cmd.AddCommand(performanceCmd)
	cmd.AddCommand(clusterCmd)
	cmd.Execute()
}
