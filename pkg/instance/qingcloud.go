package instance

import (
	"fmt"
	"time"

	"github.com/magicsong/yunify-k8s/pkg/retry"
	"github.com/yunify/qingcloud-sdk-go/client"
	"github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog/klogr"
)

const (
	DefaultCreateInstanceWait = time.Minute
	DefaultRetryCount         = 3

	ClusterNamePrefix = "YUNIFY_K8S_APP"
)

func GeneateName(clusterName string, role byte) string {
	roleName := "master"
	if role == RoleNode {
		roleName = "node"
	}
	return fmt.Sprintf("%s-%s-%s", ClusterNamePrefix, clusterName, roleName)
}

var log = klogr.New().WithName("Instance")

var _ Interface = &qingcloudInstance{}

func NewQingCloudInstanceService(instance *service.InstanceService, job *service.JobService) Interface {
	return &qingcloudInstance{
		jobService:      job,
		instanceService: instance,
	}
}

type qingcloudInstance struct {
	jobService      *service.JobService
	instanceService *service.InstanceService
}

func (q *qingcloudInstance) CreateInstances(opt *CreateInstancesOption) ([]*Instance, error) {
	input := &service.RunInstancesInput{
		Count:         &opt.Count,
		InstanceClass: &opt.InstanceClass,
		LoginKeyPair:  &opt.SSHKeyID,
		VxNets:        []*string{&opt.VxNet},
		LoginMode:     service.String("keypair"),
		InstanceName:  service.String(GeneateName(opt.Name, opt.Role)),
	}
	if opt.Role == RoleMaster {
		input.CPU = &opt.MasterCPU
		input.Memory = &opt.MasterMemory
		input.ImageID = &opt.MasterImageID
	} else if opt.Role == RoleNode {
		input.CPU = &opt.NodeCPU
		input.Memory = &opt.NodeMemory
		input.ImageID = &opt.NodeImageID
	}

	output, err := q.instanceService.RunInstances(input)
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		err := fmt.Errorf("Error in creating instances, err: %s", *output.Message)
		return nil, err
	}
	log.Info("Waiting for instance starting")
	err = client.WaitJob(q.jobService, *output.JobID, DefaultCreateInstanceWait, time.Second*5)
	if err != nil {
		return nil, err
	}
	log.Info("Machines starting successfully")
	log.Info("Waiting for instance getting its ip")
	result := make([]*Instance, 0)
	for _, i := range output.Instances {
		ins, err := client.WaitInstanceNetwork(q.instanceService, *i, DefaultCreateInstanceWait, time.Second*5)
		if err != nil {
			return nil, err
		}
		result = append(result, &Instance{
			ID: *ins.InstanceID,
			IP: *ins.VxNets[0].PrivateIP,
		})
	}
	return result, nil
}

func (q *qingcloudInstance) GetInstance(id string) (*Instance, error) {
	result, err := q.getInstancesWithRetry([]*string{&id}, DefaultRetryCount)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (q *qingcloudInstance) getInstancesWithRetry(ids []*string, retryTimes int) ([]*Instance, error) {
	input := &service.DescribeInstancesInput{
		Instances: ids,
		Verbose:   service.Int(1),
	}
	result := make([]*Instance, 0)
	err := retry.Do(retryTimes, time.Second, func() error {
		output, err := q.instanceService.DescribeInstances(input)
		if err != nil {
			log.Error(err, "error in getting instances, retry again")
			return err
		}
		if *output.RetCode != 0 {
			err := fmt.Errorf("err: %s", *output.Message)
			log.Error(err, "error in getting instances, retry again")
			return err
		}
		for _, i := range output.InstanceSet {
			result = append(result, &Instance{
				ID: *i.InstanceID,
				IP: *i.VxNets[0].PrivateIP,
			})
		}
		return nil
	})
	return result, err
}

func (q *qingcloudInstance) DeleteInstance(instanceID string) error {
	return nil
}
