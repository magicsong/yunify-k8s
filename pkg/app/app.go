package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	accesskey "github.com/magicsong/yunify-k8s/pkg/access-key"
	"github.com/magicsong/yunify-k8s/pkg/api"
	"github.com/magicsong/yunify-k8s/pkg/instance"
	"github.com/magicsong/yunify-k8s/pkg/ssh"
	"github.com/magicsong/yunify-k8s/pkg/sshkey"
	"github.com/magicsong/yunify-k8s/pkg/tag"
	"k8s.io/klog"
)

type App interface {
	RunCreate(*api.CreateClusterOption) error
	RunDelete(*api.DeleteClusterOption) error
}

func NewApp(configFile string) App {
	return &app{
		configFile: configFile,
	}
}

type app struct {
	instanceIface instance.Interface
	sshKeyIface   sshkey.Interface
	tagService    tag.Interface
	configFile    string
}

func tagName(name string) string {
	return fmt.Sprintf("K8S-Cluster-%s", name)
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
	keyHelper := accesskey.NewQingCloudAccessKeyHelper(zone, a.configFile)
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
	tagService, _ := qcService.Tag(zone)
	a.tagService = tag.NewQingCloudTagService(tagService, userid)
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
	klog.Info("Prepare Tag")
	tag := tagName(opt.ClusterName)
	id, err := a.tagService.GetTagClusterByName(tag)
	if err != nil {
		klog.Error("Failed to get current tag")
		return err
	}
	var tagID string
	if id != nil {
		tagID = id.TagID
	} else {
		tagID, err = a.tagService.CreateTag(tag)
		if err != nil {
			klog.Errorf("Failed to create tag %s", tag)
			return err
		}
	}
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
	machines := []string{}
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
		machines = append(machines, master.ID)
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
			machines = append(machines, machine.ID)
			nodes = append(nodes, machine)
			klog.Infof("Nodes creating done, id=%s, ip=%s", machine.ID, machine.IP)
		}
	}()

	klog.Infoln("Waiting for machines to start")
	wg.Wait()
	if len(errs) != 0 {
		return fmt.Errorf("Creating Machines failed, errs: %+v", errs)
	}

	klog.Infoln("Tagging all machines")
	err = a.tagService.TagInstances(tagID, machines)
	if err != nil {
		return err
	}
	klog.Infoln("Machines are ready, bring the cluster up")
	_, err = bootstrapMaster(master, opt)
	if err != nil {
		klog.Errorln("Failed to bootstrap master node")
		return err
	}
	klog.Info("Apply CNI")
	return nil
}

func generateKubeadmInitCmd(opt api.NetworkOption, version string) (string, error) {
	if opt.PodNetWorkCIDR == "" {
		return "", fmt.Errorf("Must specify a network for pod")
	}

	if opt.CNIName == api.CalicoCNI || opt.CNIName == api.FlannelCNI {
		return fmt.Sprintf("kubeadm init --pod-network-cidr=%s --kubernetes-version=v%s", opt.PodNetWorkCIDR, version), nil
	}

	return "", fmt.Errorf("CNI plugin %s is not supported right now", opt.CNIName)
}

func bootstrapMaster(master *instance.Instance, opt *api.CreateClusterOption) (string, error) {
	cmd, err := generateKubeadmInitCmd(opt.NetworkOption, opt.KubernetesVersion)
	if err != nil {
		return "", err
	}
	output, err := ssh.QuickConnectAndRun(master.IP, "swapoff -a;"+cmd)
	defer klog.Infoln(string(output))
	if err != nil {
		klog.Errorln("Failed to run 'kubeadm init'")
		return "", err
	}
	klog.Info("Getting 'kubeadm join'")
	return GetKubeJoinFromOutput(string(output)), nil
}

func buildShellScripts(scripts []string) string {
	var buf bytes.Buffer
	buf.WriteString("#!/bin/bash\n")
	buf.WriteString("swapoff -a\n")
	for _, s := range scripts {
		buf.WriteString(s)
		buf.WriteString("\n")
	}
	return buf.String()
}

func GetKubeJoinFromOutput(output string) string {
	output = strings.TrimSpace(output)
	index := strings.LastIndex(output, "kubeadm join")
	return output[index:]
}
