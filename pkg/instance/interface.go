package instance

import (
	"github.com/magicsong/yunify-k8s/pkg/api"
)

type Instance struct {
	ID string
	IP string
}

type CreateInstancesOption struct {
	Name          string
	VxNet         string
	SSHKeyID      string
	Count         int
	Role          byte
	InstanceClass int
	api.ImagesPreset
}

type Interface interface {
	CreateInstances(*CreateInstancesOption) ([]*Instance, error)
	DeleteInstances(instanceID []string) error
	GetInstance(string) (*Instance, error)
	StopInstances(...string) error
}
