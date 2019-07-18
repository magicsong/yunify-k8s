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
	createCmd.AddCommand(createClusterCmd)
	createClusterOpt = new(api.CreateClusterOption)
	createClusterCmd.Flags().StringVarP(&createClusterOpt.KubernetesVersion, "k8sVersion", "k", "1.13.1", "specify k8s version of cluster")
	createClusterCmd.Flags().StringVarP(&createClusterOpt.PodNetWorkCIDR, "podcidr", "p", "10.233.0.0/16", "specify PodNetWorkCIDR")
	createClusterCmd.Flags().IntVarP(&createClusterOpt.NodeCount, "nodecount", "c", 2, "specify the number of nodes")
	createClusterCmd.Flags().StringVar(&createClusterOpt.CNIName, "cni", "calico", "cni plugin to use")
	createClusterCmd.Flags().IntVar(&createClusterOpt.InstanceClass, "class", 101, "instance class of machine,available values: 0, 1, 2, 3, 4, 5, 6, 100, 101, 200, 201, 300, 301")
	createClusterCmd.Flags().BoolVarP(&createClusterOpt.ScpKubeConfigToLocal, "scp-kubeconfig", "s", false, "specify whether copy kubeconfig to local")
	createClusterCmd.Flags().StringVar(&createClusterOpt.LocalKubeConfigPath, "kubeconfig-path", ".", "specify the path where kubeconfig copy to")
	createCmd.MarkFlagRequired("vxnet")
}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "create a new cluster",
	Long: `create a new cluster, for example:
  qks create my-k8s-cluster --vxnet=vxxxxx`,
	ValidArgs: []string{"clusterName"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createClusterOpt.ClusterName = args[0]
		createClusterOpt.Zone = zone
		createClusterOpt.VxNet = vxnet
		createClusterOpt.UseExistKey = useExistKey
		toRun := app.NewApp(cfgFile)
		err := toRun.RunCreate(createClusterOpt)
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}
	},
}
