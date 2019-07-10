package tag

type TagCluster struct {
	TagID     string
	Instances []string
}

type Interface interface {
	CreateTag(string) (string, error)
	DeleteTag(string) error
	GetTagClusterByName(string) (*TagCluster, error)
	TagInstances(string, []string) error
}
