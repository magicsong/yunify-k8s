package cmd

import (
	"os"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/app"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var createClusterOpt *api.CreateClusterOption

func init() {
	rootCmd.AddCommand(createCmd)
	createClusterOpt = new(api.CreateClusterOption)
	createCmd.Flags().StringVarP(&createClusterOpt.KubernetesVersion, "k8sVersion", "k", "13.1", "specify k8s version of cluster")
	createCmd.Flags().StringVarP(&createClusterOpt.PodNetWorkCIDR, "podcidr", "p", "10.10.0.0/16", "specify PodNetWorkCIDR")
	createCmd.Flags().IntVarP(&createClusterOpt.NodeCount, "nodecount", "c", 2, "specify the number of nodes")
	createCmd.Flags().StringVarP(&createClusterOpt.VxNet, "vxnet", "x", "", "specify the vxnet")
	createCmd.Flags().StringVar(&createClusterOpt.CNIName, "cni", "calico", "cni plugin to use")
	createCmd.Flags().IntVar(&createClusterOpt.InstanceClass, "class", 101, "instance class of machine,available values: 0, 1, 2, 3, 4, 5, 6, 100, 101, 200, 201, 300, 301")
	createCmd.Flags().BoolVar(&createClusterOpt.UseExistKey, "use-old-key", false, "specify whether create or reuse former ssh key to connect cluster")
	createCmd.MarkFlagRequired("vxnet")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new cluster",
	Long: `create a new cluster, for example:
  qks create my-k8s-cluster --vxnet=vxxxxx`,
	ValidArgs: []string{"clusterName"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createClusterOpt.ClusterName = args[0]
		createClusterOpt.Zone = zone
		toRun := app.NewApp(cfgFile)
		err := toRun.RunCreate(createClusterOpt)
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}
	},
}
