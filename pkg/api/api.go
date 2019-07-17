package api

const (
	ErrorK8sVersionNotSupport = "Currently we do not support k8s version %s"
	SSHKeyName                = "DO_NOT_REMOVE_K8S_KEY"
	CalicoCNI                 = "calico"
	FlannelCNI                = "flannel"
)

const (
	RoleMaster byte = iota
	RoleNode
)

type CreateClusterOption struct {
	ClusterName       string
	KubernetesVersion string
	NodeCount         int
	VxNet             string
	InstanceClass     int
	Zone              string
	NetworkOption
	UseExistKey          bool
	ScpKubeConfigToLocal bool
	LocalKubeConfigPath  string
}

type NetworkOption struct {
	CNIName        string
	PodNetWorkCIDR string
}

type DeleteClusterOption struct {
	ClusterName string
	ForceDelete bool
	Zone        string
}

type CreateImageOption struct {
	BaseImage       string
	ManifestFolders []string
	Scripts         []string
	ImageName       string
	Role            byte
	VxNet           string
	UseExistKey     bool
	Zone            string
	DeleteMachine   bool
}
