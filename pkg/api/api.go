package api

const (
	ErrorK8sVersionNotSupport = "Currently we do not support k8s version %s"
	SSHKeyName                = "DO_NOT_REMOVE_K8S_KEY"
	CalicoCNI                 = "calico"
	FlannelCNI                = "flannel"
)

type CreateClusterOption struct {
	KubernetesVersion string
	NodeCount         int
	VxNet             string
	InstanceClass     int
	Zone              string
	NetworkOption
	UseExistKey bool
}

type NetworkOption struct {
	CNIName        string
	PodNetWorkCIDR string
}
