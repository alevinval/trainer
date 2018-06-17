package main

import (
	"github.com/spf13/cobra"
)

var (
	prefix              string
	enrichStravaCsvPath string
	filterByName  	    string
	cmd                 *cobra.Command
)

func init() {
	cmd = &cobra.Command{Use: "trainer"}
	cmd.PersistentFlags().StringVar(&prefix, "prefix", "", "only process files matching the prefix")
	cmd.PersistentFlags().StringVar(&enrichStravaCsvPath, "strava-csv-enrich", "", "enrich metadata from a csv file")
	cmd.PersistentFlags().StringVar(&filterByName, "filter-by-name", "", "Ignores activities whose name does not match the filter")
}

func main() {
	cmd.AddCommand(performanceCmd)
	cmd.AddCommand(clusterCmd)
	cmd.Execute()
}
