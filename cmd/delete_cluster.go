package cmd

import (
	"os"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/app"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var deleteClusterOpt *api.DeleteClusterOption

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)
	deleteClusterOpt = new(api.DeleteClusterOption)
}

var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "delete a new cluster",
	Long: `delete a new cluster, for example:
  qks delete my-k8s-cluster`,
	ValidArgs: []string{"clusterName"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteClusterOpt.ClusterName = args[0]
		deleteClusterOpt.Zone = zone
		toRun := app.NewApp(cfgFile)
		err := toRun.RunDelete(deleteClusterOpt)
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}
	},
}
