package main

import (
	"github.com/spf13/cobra"
)

var (
	prefix              string
	enrichStravaCsvPath string
	cmd                 *cobra.Command
)

func init() {
	cmd = &cobra.Command{Use: "trainer"}
	cmd.PersistentFlags().StringVar(&prefix, "prefix", "", "only process files matching the prefix")
	cmd.PersistentFlags().StringVar(&enrichStravaCsvPath, "strava-csv-enrich", "", "enrich metadata from a csv file")
}

func main() {
	cmd.AddCommand(performanceCmd)
	cmd.AddCommand(clusterCmd)
	cmd.Execute()
}
