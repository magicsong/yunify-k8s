package app

import (
	"fmt"
	"time"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"k8s.io/klog"
)

func (a *app) RunDelete(opt *api.DeleteClusterOption) error {
	start := time.Now()
	defer func() {
		runningTime := time.Since(start)
		klog.Infof("Finished, time cost(s): %d", runningTime/time.Second)
	}()
	err := a.validateDeleteInput(opt)
	if err != nil {
		return err
	}
	err = a.init(opt.Zone)
	if err != nil {
		klog.Error("Falied to init command")
		return err
	}
	return a.runDelete(opt)
}

func (a *app) validateDeleteInput(opt *api.DeleteClusterOption) error {
	if opt.ClusterName == "" {
		return fmt.Errorf("ClusterName cannot be empty")
	}
	return nil
}

func (a *app) runDelete(opt *api.DeleteClusterOption) error {
	tagInstances, err := a.tagService.GetTagClusterByName(tagName(opt.ClusterName))
	if err != nil {
		klog.Errorf("Failed to get instances of cluster %s", opt.ClusterName)
		return err
	}
	if tagInstances == nil {
		err = fmt.Errorf("Cannot find the cluster %s in zone %s", opt.ClusterName, opt.Zone)
		return err
	}
	klog.Info("Begin to terminate cluster machines")
	err = a.instanceIface.DeleteInstances(tagInstances.Instances)
	if err != nil {
		return err
	}

	klog.Info("Deleting tag")
	err = a.tagService.DeleteTag(tagInstances.TagID)
	if err != nil {
		return err
	}
	klog.Info("Cluster has been successfully deleted")
	return nil
}
