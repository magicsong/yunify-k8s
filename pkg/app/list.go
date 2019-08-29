package app

import (
	"github.com/magicsong/yunify-k8s/pkg/api"
	"k8s.io/klog"
)

func (a *app) getClusters() error {
	tags, err := a.tagService.GetTags(api.ClusterTagPrefix)
	if err != nil {
		klog.Errorln("Failed to get tags")
	}
	for _, t := range tags {
		klog.Infof("Get cluster [%s]", t[len(api.ClusterTagPrefix):])
	}
	return nil
}

func (a *app) RunList(zone string) error {
	err := a.init(zone)
	if err != nil {
		klog.Error("Falied to init command")
		return err
	}
	return a.getClusters()
}
