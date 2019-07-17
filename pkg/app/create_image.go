package app

import (
	"fmt"
	"os/exec"
	"path"
	"time"

	"github.com/magicsong/yunify-k8s/pkg/instance"

	"github.com/magicsong/yunify-k8s/pkg/api"
	"k8s.io/klog"
)

var defaultImage = api.ImagesPreset{
	NodeImageID:   "img-rfubqmqn",
	MasterImageID: "img-ybttnmjg",
	NodeCPU:       4,
	NodeMemory:    4096,
	MasterCPU:     4,
	MasterMemory:  4096,
	CNIYamlPath:   "/root/CNI",
}

func prepareLocalSSHBeforeTransfering(masterip string) error {
	removeIP := exec.Command("ssh-keygen", "-R", masterip)
	bytes, err := removeIP.CombinedOutput()
	klog.V(2).Info(string(bytes))
	if err != nil {
		klog.Warningf("Failed to remove %s in known_hosts", masterip)
	}
	addHost := exec.Command("bash", "-c", "ssh-keyscan -H "+masterip+" >>~/.ssh/known_hosts")
	bytes, err = addHost.CombinedOutput()
	klog.V(2).Info(string(bytes))
	return err
}

func (a *app) createImageInstance(opt *api.CreateImageOption, sshkey string) (*instance.Instance, error) {
	createInstanceOpt := &instance.CreateInstancesOption{
		Name:          "ImageBuilder_" + opt.ImageName,
		VxNet:         opt.VxNet,
		SSHKeyID:      sshkey,
		Count:         1,
		Role:          opt.Role,
		InstanceClass: 101,
		ImagesPreset:  defaultImage,
	}
	if opt.Role == api.RoleMaster {
		createInstanceOpt.MasterImageID = opt.BaseImage
	} else {
		createInstanceOpt.NodeImageID = opt.BaseImage
	}
	instances, err := a.instanceIface.CreateInstances(createInstanceOpt)
	if err != nil {
		return nil, err
	}
	return instances[0], nil
}

func (a *app) RunCreateImage(opt *api.CreateImageOption) error {
	start := time.Now()
	defer func() {
		runningTime := time.Since(start)
		klog.Infof("Finished, time cost(s): %d", runningTime/time.Second)
	}()
	err := a.init(opt.Zone)
	if err != nil {
		klog.Error("Falied to init command")
		return err
	}
	return a.runCreateImage(opt)
}

func (a *app) runCreateImage(opt *api.CreateImageOption) error {
	klog.Info("Prepare ssh key")
	keyid, err := a.prepareSSHKey(opt.UseExistKey)
	if err != nil {
		return err
	}
	klog.Info("Creating machine to bulid image")
	inst, err := a.createImageInstance(opt, keyid)
	if err != nil {
		klog.Error("Failed to create instance")
		return err
	}
	klog.Infof("instance %s [%s] is up ,begin to run image scripts", inst.ID, inst.IP)
	klog.Infof("Add %s to local known_hosts", inst.IP)
	err = prepareLocalSSHBeforeTransfering(inst.IP)
	if err != nil {
		klog.Error("Failed to add host to known_hosts")
		return err
	}
	klog.Info("Transfer files")
	for _, folder := range opt.ManifestFolders {
		err = transferFolder(inst.IP, folder)
		if err != nil {
			klog.Errorf("Failed to scp folder %s to remote", folder)
			return err
		}
	}
	for _, file := range opt.Scripts {
		err = transferFile(inst.IP, file)
		if err != nil {
			klog.Errorf("Failed to scp file %s to remote", file)
			return err
		}
	}
	klog.Infof("Transfer done")
	if opt.DeleteMachine {
		klog.Infof("Begin to tear down machine %s", inst.IP)
		err = a.instanceIface.DeleteInstances([]string{inst.ID})
		if err != nil {
			klog.Warningf("Failed to delete machine %s, you have to  do it manually. Err: %s", inst.ID, err.Error())
		}
	}
	return nil
}

func transferFolder(ip, folder string) error {
	//scp -r ../vm-scripts root@$ip:/root/vm-scripts
	p := path.Base(folder)
	cmd := exec.Command("scp", "-r", folder, fmt.Sprintf("root@%s:/root/%s", ip, p))
	bytes, err := cmd.CombinedOutput()
	klog.V(2).Info(string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func transferFile(ip, filePath string) error {
	p := path.Base(filePath)
	cmd := exec.Command("scp", filePath, fmt.Sprintf("root@%s:/tmp/%s", ip, p))
	bytes, err := cmd.CombinedOutput()
	klog.V(2).Info(string(bytes))
	if err != nil {
		return err
	}
	return nil
}
