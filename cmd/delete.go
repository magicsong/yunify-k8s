package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete clusters or vm images",
	Long: `qks is a CLI to rapidly create/detele a kubernetes in qingcloud. for example:
  qks create my-k8s-cluster --vxnet=vxxxxx`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
