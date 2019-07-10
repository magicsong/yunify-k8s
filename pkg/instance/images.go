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
}

func init() {
	PresetKubernetes = make(map[string]ImagesPreset)
	PresetKubernetes["12.4"] = ImagesPreset{
		KubernetesVersion: "12.4",
		NodeImageID:       "img-rfubqmqn",
		MasterImageID:     "img-ybttnmjg",
		NodeCPU:           4,
		NodeMemory:        4096,
		MasterCPU:         4,
		MasterMemory:      4096,
	}
}
