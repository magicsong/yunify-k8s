package cmd

import (
	"io/ioutil"
	"os"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/app"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

var createImageOpt *api.CreateImageOption
var createImageYaml string

func init() {
	createCmd.AddCommand(createImageCmd)
	createImageOpt = new(api.CreateImageOption)
	createImageCmd.Flags().StringVarP(&createImageYaml, "yaml", "Y", "", "Use yaml instead of Command line")
	createImageCmd.Flags().StringVarP(&createImageOpt.InstanceInfo.BaseImage, "base-image", "i", "xenial5x64a", "specify the base image to work on")
	createImageCmd.Flags().BoolVarP(&createImageOpt.DeleteMachine, "delete-machine", "D", false, "specify whether deleting  machine or not in the end")
	createImageCmd.Flags().StringArrayVarP(&createImageOpt.Manifest.Folders, "scripts-folder", "F", nil, "folders will be upload to image")
	createImageCmd.Flags().StringArrayVarP(&createImageOpt.Manifest.Scripts, "script-file", "f", nil, "files to be uploaded")
	createImageCmd.Flags().StringVarP(&createImageOpt.EntryPoint, "entry", "e", "", "The script to run in the image")
}

var createImageCmd = &cobra.Command{
	Use:   "image",
	Short: "delete a new cluster",
	Long: `delete a new cluster, for example:
  qks create image test-image -x=vxnetxxx`,
	Run: func(cmd *cobra.Command, args []string) {
		if createImageYaml != "" {
			bytes, err := ioutil.ReadFile(createImageYaml)
			if err != nil {
				klog.Errorf("Failed to read yaml,err: %s", err.Error())
				os.Exit(1)
			}
			err = yaml.Unmarshal(bytes, createImageOpt)
			if err != nil {
				klog.Errorf("Failed to parse yaml,err: %s", err.Error())
				os.Exit(1)
			}
		} else {
			if len(args) != 1 {
				klog.Error("Must specify a image name, for example 'qks create image test-image'")
				os.Exit(1)
			}
			createImageOpt.ImageName = args[0]
			createImageOpt.InstanceInfo.Zone = zone
			createImageOpt.InstanceInfo.VxNet = vxnet
			createImageOpt.InstanceInfo.UseExistKey = useExistKey
		}
		toRun := app.NewApp(cfgFile)
		err := toRun.RunCreateImage(createImageOpt)
		if err != nil {
			klog.Errorln(err)
			os.Exit(1)
		}
	},
}
