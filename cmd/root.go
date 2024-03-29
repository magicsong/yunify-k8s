/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var cfgFile string
var zone string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qks",
	Short: "qks is a CLI to rapidly create/detele a kubernetes in qingcloud",
	Long: `qks is a CLI to rapidly create/detele a kubernetes in qingcloud. for example:
  qks create my-k8s-cluster --vxnet=vxxxxx`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	klog.InitFlags(nil)
	goflag.Set("logtostderr", "true")
	goflag.Set("alsologtostderr", "false")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.qingcloud/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&zone, "zone", "z", "ap2a", "specify zone to delete cluster")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)
}
