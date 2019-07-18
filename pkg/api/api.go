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
	ImageName     string              `yaml:"name,omitempty"`
	Manifest      CreateImageManifest `yaml:"manifest,omitempty"`
	DeleteMachine bool                `yaml:"deleteMachine,omitempty"`
	EntryPoint    string              `yaml:"entryPoint,omitempty"`
	InstanceInfo  InstanceInfo        `yaml:"instanceInfo,omitempty"`
}
type CreateImageManifest struct {
	Folders []string `yaml:"manifestFolders,omitempty,flow"`
	Scripts []string `yaml:"scripts,omitempty,flow"`
}
type InstanceInfo struct {
	BaseImage   string `yaml:"baseImage,omitempty"`
	Role        byte   `yaml:"role,omitempty"`
	VxNet       string `yaml:"vxNet,omitempty"`
	UseExistKey bool   `yaml:"useExistKey,omitempty"`
	Zone        string `yaml:"zone,omitempty"`
}
