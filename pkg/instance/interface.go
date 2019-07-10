package instance

type Instance struct {
	ID string
	IP string
}

const (
	RoleMaster byte = iota
	RoleNode
)

type CreateInstancesOption struct {
	VxNet         string
	SSHKeyID      string
	Count         int
	Role          byte
	InstanceClass int
	ImagesPreset
}

type Interface interface {
	CreateInstances(*CreateInstancesOption) ([]*Instance, error)
	DeleteInstance(instanceID string) error
	GetInstance(string) (*Instance, error)
}
