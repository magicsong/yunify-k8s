package cmd

import (
	"github.com/spf13/cobra"
)

var vxnet string
var useExistKey bool

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create clusters or vm images",
	Long: `qks is a CLI to rapidly create/detele a kubernetes in qingcloud. for example:
  qks create my-k8s-cluster --vxnet=vxxxxx`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().StringVarP(&vxnet, "vxnet", "x", "", "specify the vxnet")
	createCmd.PersistentFlags().BoolVar(&useExistKey, "use-old-key", true, "specify whether create or reuse former ssh key to connect machines")
}
