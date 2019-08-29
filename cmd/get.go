package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get clusters or vm images",
	Long: `qks is a CLI to rapidly get/detele a kubernetes in qingcloud. for example:
  qks get cluster`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
