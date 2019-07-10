package app

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	accesskey "github.com/magicsong/yunify-k8s/pkg/access-key"
	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/instance"
	"github.com/magicsong/yunify-k8s/pkg/ssh"
	"github.com/magicsong/yunify-k8s/pkg/sshkey"
	"k8s.io/klog"
)

type App interface {
	RunCreate(opt *api.CreateClusterOption) error
	RunDelete() error
}

func NewApp() App {
	return &app{}
}

type app struct {
	instanceIface instance.Interface
	sshKeyIface   sshkey.Interface
}

func (a *app) RunDelete() error {
	return nil
}

func (a *app) validateCreateInput(opt *api.CreateClusterOption) error {
	if opt.ClusterName == "" {
		return fmt.Errorf("ClusterName cannot be empty")
	}
	return nil
}
func (a *app) RunCreate(opt *api.CreateClusterOption) error {
	start := time.Now()
	defer func() {
		runningTime := time.Since(start)
		klog.Infof("Finished, time cost(s): %d", runningTime/time.Second)
	}()
	err := a.validateCreateInput(opt)
	if err != nil {
		return err
	}
	err = a.init(opt.Zone)
	if err != nil {
		klog.Error("Falied to init command")
		return err
	}
	return a.runCreate(opt)
}

func (a *app) init(zone string) error {
	klog.Info("Init qingcloud service")
	keyHelper := accesskey.NewQingCloudAccessKeyHelper(zone, "")
	err := keyHelper.Init()
	if err != nil {
		return err
	}
	userid := keyHelper.GetUserID()
	qcService := keyHelper.GetService()
	instanceService, _ := qcService.Instance(zone)
	jobService, _ := qcService.Job(zone)
	a.instanceIface = instance.NewQingCloudInstanceService(instanceService, jobService)
	keyService, _ := qcService.KeyPair(zone)
	a.sshKeyIface = sshkey.NewQingCloudKeyPairService(keyService, userid)
	return nil
}

func (a *app) prepareSSHKey(useExistKey bool) (string, error) {
	output, err := ioutil.ReadFile(ssh.GetDefaultPublicKeyFile())
	if err != nil {
		klog.Errorln("Failed to read ssh public key")
		return "", err
	}
	if useExistKey {
		klog.Info("Try to get exsit keypair")
		key, err := a.sshKeyIface.GetKeyPairByName(api.SSHKeyName)
		if err != nil {
			return "", err
		}
		if key != "" {
			return key, err
		}
		klog.Warning("Cannot find any exist key, will create a new one")
	}
	klog.Info("Try to create a new ssh key")
	return a.sshKeyIface.CreateSSHKey(api.SSHKeyName, string(output))
}

func (a *app) runCreate(opt *api.CreateClusterOption) error {
	klog.Info("Prepare ssh key")
	keyid, err := a.prepareSSHKey(opt.UseExistKey)
	if err != nil {
		return err
	}
	//create master
	var wg sync.WaitGroup
	klog.Infoln("Creating Master")
	if _, ok := instance.PresetKubernetes[opt.KubernetesVersion]; !ok {
		return fmt.Errorf(api.ErrorK8sVersionNotSupport, opt.KubernetesVersion)
	}
	wg.Add(1)
	var master *instance.Instance
	errs := make([]error, 0)
	go func() {
		defer wg.Done()
		createMasterOpt := &instance.CreateInstancesOption{
			Name:          opt.ClusterName,
			VxNet:         opt.VxNet,
			Count:         1,
			Role:          instance.RoleMaster,
			ImagesPreset:  instance.PresetKubernetes[opt.KubernetesVersion],
			InstanceClass: opt.InstanceClass,
			SSHKeyID:      keyid,
		}
		instances, err := a.instanceIface.CreateInstances(createMasterOpt)
		if err != nil {
			errs = append(errs, err)
			return
		}
		master = instances[0]
		klog.Infof("Master creating done, id=%s, ip=%s", master.ID, master.IP)
	}()
	//creating nodes
	wg.Add(1)
	var nodes []*instance.Instance

	go func() {
		defer wg.Done()
		createNodesOpt := &instance.CreateInstancesOption{
			Name:          opt.ClusterName,
			VxNet:         opt.VxNet,
			Count:         opt.NodeCount,
			Role:          instance.RoleNode,
			ImagesPreset:  instance.PresetKubernetes[opt.KubernetesVersion],
			InstanceClass: opt.InstanceClass,
			SSHKeyID:      keyid,
		}
		instances, err := a.instanceIface.CreateInstances(createNodesOpt)
		if err != nil {
			errs = append(errs, err)
			return
		}
		for _, machine := range instances {
			nodes = append(nodes, machine)
			klog.Infof("Nodes creating done, id=%s, ip=%s", machine.ID, machine.IP)
		}
	}()

	klog.Infoln("Waiting for machines to start")
	wg.Wait()
	if len(errs) != 0 {
		return fmt.Errorf("Creating Machines failed, errs: %+v", errs)
	}
	klog.Infoln("Machines are ready, bring the cluster up")
	_, err = bootstrapMaster(master, opt.NetworkOption)
	if err != nil {
		klog.Errorln("Failed to bootstrap master node")
		return err
	}
	return nil
}

func generateKubeadmInitCmd(opt api.NetworkOption) (string, error) {
	if opt.PodNetWorkCIDR == "" {
		return "", fmt.Errorf("Must specify a network for pod")
	}

	if opt.CNIName == api.CalicoCNI || opt.CNIName == api.FlannelCNI {
		return fmt.Sprintf("kubeadm init --pod-network-cidr=%s", opt.PodNetWorkCIDR), nil
	}

	return "", fmt.Errorf("CNI plugin %s is not supported right now", opt.CNIName)
}

func bootstrapMaster(master *instance.Instance, opt api.NetworkOption) (string, error) {
	cmd, err := generateKubeadmInitCmd(opt)
	if err != nil {
		return "", err
	}
	output, err := ssh.QuickConnectAndRun(master.IP, cmd)
	defer klog.Infoln(string(output))
	if err != nil {
		klog.Errorln("Failed to run 'kubeadm init'")
		return "", err
	}
	return "", nil
}
