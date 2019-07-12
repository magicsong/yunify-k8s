package instance

var PresetKubernetes map[string]ImagesPreset

type ImagesPreset struct {
	KubernetesVersion string
	NodeImageID       string
	MasterImageID     string
	NodeCPU           int
	NodeMemory        int
	MasterCPU         int
	MasterMemory      int
	CNIYamlPath       string
}

func init() {
	PresetKubernetes = make(map[string]ImagesPreset)
	PresetKubernetes["1.13.1"] = ImagesPreset{
		KubernetesVersion: "1.13.1",
		NodeImageID:       "img-rfubqmqn",
		MasterImageID:     "img-ybttnmjg",
		NodeCPU:           4,
		NodeMemory:        4096,
		MasterCPU:         4,
		MasterMemory:      4096,
		CNIYamlPath:       "/root/CNI",
	}
}
