package cmd

import (
	"os"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/app"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var createImageOpt *api.CreateImageOption

func init() {
	createCmd.AddCommand(createImageCmd)
	createImageOpt = new(api.CreateImageOption)
	createImageCmd.Flags().StringVarP(&createImageOpt.BaseImage, "base-image", "i", "", "specify the base image to work on")
	createImageCmd.Flags().BoolVarP(&createImageOpt.DeleteMachine, "delete-machine", "D", false, "specify whether deleting  machine or not in the end")
	createImageCmd.Flags().StringArrayVarP(&createImageOpt.ManifestFolders, "scripts-folder", "F", nil, "folders will be upload to image")
}

var createImageCmd = &cobra.Command{
	Use:   "image",
	Short: "delete a new cluster",
	Long: `delete a new cluster, for example:
  qks create image test-image -x=vxnetxxx`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createImageOpt.ImageName = args[0]
		createImageOpt.Zone = zone
		createImageOpt.VxNet = vxnet
		createImageOpt.UseExistKey = useExistKey
		toRun := app.NewApp(cfgFile)
		err := toRun.RunCreateImage(createImageOpt)
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}
	},
}
