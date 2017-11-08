package main

import (
	"github.com/spf13/cobra"
)

var (
	prefix string
	cmd    *cobra.Command
)

func init() {
	cmd = &cobra.Command{Use: "trainer"}
	cmd.PersistentFlags().StringVar(&prefix, "prefix", "", "only process files matching the prefix")
}

func main() {
	cmd.AddCommand(performanceCmd)
	cmd.AddCommand(clusterCmd)
	cmd.Execute()
}
