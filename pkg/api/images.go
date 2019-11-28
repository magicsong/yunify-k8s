package api

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
	CNICmd            string
}

func init() {
	PresetKubernetes = make(map[string]ImagesPreset)
	PresetKubernetes["1.13.1"] = ImagesPreset{
		KubernetesVersion: "1.13.1",
		NodeImageID:       "img-sykyoovw",
		MasterImageID:     "img-ybttnmjg",
		NodeCPU:           4,
		NodeMemory:        4096,
		MasterCPU:         4,
		MasterMemory:      4096,
		CNIYamlPath:       "/root/CNI",
		CNICmd:            "cni.sh",
	}

	PresetKubernetes["1.15.2"] = ImagesPreset{
		KubernetesVersion: "1.15.2",
		NodeImageID:       "img-sykyoovw",
		MasterImageID:     "img-kj5hg0fe",
		NodeCPU:           4,
		NodeMemory:        4096,
		MasterCPU:         4,
		MasterMemory:      4096,
		CNIYamlPath:       "/root/CNI",
		CNICmd:            "cni.sh",
	}
	PresetKubernetes["1.15.5"] = ImagesPreset{
		KubernetesVersion: "1.15.5",
		NodeImageID:       "img-sykyoovw",
		MasterImageID:     "img-kj5hg0fe",
		NodeCPU:           2,
		NodeMemory:        2048,
		MasterCPU:         2,
		MasterMemory:      4096,
		CNIYamlPath:       "/root/CNI",
		CNICmd:            "cni.sh",
	}
}
