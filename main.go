package main

import (
	"flag"
	"os"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/app"
	"k8s.io/klog"
)

var createClusterOpt *api.CreateClusterOption

func init() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Set("alsologtostderr", "false")

	createClusterOpt = new(api.CreateClusterOption)
	flag.StringVar(&createClusterOpt.ClusterName, "name", "", "The name of cluster, must unique in a zone, required")
	flag.StringVar(&createClusterOpt.KubernetesVersion, "k", "13.1", "specify k8s version of cluster")
	flag.StringVar(&createClusterOpt.PodNetWorkCIDR, "n", "10.10.0.0/16", "specify PodNetWorkCIDR")
	flag.IntVar(&createClusterOpt.NodeCount, "c", 2, "specify the number of nodes")
	flag.StringVar(&createClusterOpt.VxNet, "vxnet", "", "specify the vxnet, required")
	flag.StringVar(&createClusterOpt.CNIName, "cni", "calico", "cni plugin to use")
	flag.IntVar(&createClusterOpt.InstanceClass, "class", 101, "instance class of machine,available values: 0, 1, 2, 3, 4, 5, 6, 100, 101, 200, 201, 300, 301")
	flag.StringVar(&createClusterOpt.Zone, "z", "ap2a", "specify zone to create cluster")
	flag.BoolVar(&createClusterOpt.UseExistKey, "use-old-key", false, "specify whether create or reuse former ssh key to connect cluster")
}
func main() {
	flag.Parse()
	cmd := app.NewApp()
	err := cmd.RunCreate(createClusterOpt)
	if err != nil {
		klog.Errorln(err)
		os.Exit(1)
	}
}
