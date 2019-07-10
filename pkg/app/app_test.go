package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	It("Should be able to extract kubeadm join", func() {
		case1 := `[init] Using Kubernetes version: v1.13.8
		[preflight] Running pre-flight checks
		[preflight] Pulling images required for setting up a Kubernetes cluster
		[preflight] This might take a minute or two, depending on the speed of your internet connection
		[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
		[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
		[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
		[kubelet-start] Activating the kubelet service
		[certs] Using certificateDir folder "/etc/kubernetes/pki"
		[certs] Generating "etcd/ca" certificate and key
		[certs] Generating "etcd/peer" certificate and key
		[certs] etcd/peer serving cert is signed for DNS names [i-z0333hl2 localhost] and IPs [192.168.97.4 127.0.0.1 ::1]
		[certs] Generating "etcd/healthcheck-client" certificate and key
		[certs] Generating "apiserver-etcd-client" certificate and key
		[certs] Generating "etcd/server" certificate and key
		[certs] etcd/server serving cert is signed for DNS names [i-z0333hl2 localhost] and IPs [192.168.97.4 127.0.0.1 ::1]
		[certs] Generating "ca" certificate and key
		[certs] Generating "apiserver" certificate and key
		[certs] apiserver serving cert is signed for DNS names [i-z0333hl2 kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local] and IPs [10.96.0.1 192.168.97.4]
		[certs] Generating "apiserver-kubelet-client" certificate and key
		[certs] Generating "front-proxy-ca" certificate and key
		[certs] Generating "front-proxy-client" certificate and key
		[certs] Generating "sa" key and public key
		[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
		[kubeconfig] Writing "admin.conf" kubeconfig file
		[kubeconfig] Writing "kubelet.conf" kubeconfig file
		[kubeconfig] Writing "controller-manager.conf" kubeconfig file
		[kubeconfig] Writing "scheduler.conf" kubeconfig file
		[control-plane] Using manifest folder "/etc/kubernetes/manifests"
		[control-plane] Creating static Pod manifest for "kube-apiserver"
		[control-plane] Creating static Pod manifest for "kube-controller-manager"
		[control-plane] Creating static Pod manifest for "kube-scheduler"
		[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
		[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
		[apiclient] All control plane components are healthy after 20.004137 seconds
		[uploadconfig] storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
		[kubelet] Creating a ConfigMap "kubelet-config-1.13" in namespace kube-system with the configuration for the kubelets in the cluster
		[patchnode] Uploading the CRI Socket information "/var/run/dockershim.sock" to the Node API object "i-z0333hl2" as an annotation
		[mark-control-plane] Marking the node i-z0333hl2 as control-plane by adding the label "node-role.kubernetes.io/master=''"
		[mark-control-plane] Marking the node i-z0333hl2 as control-plane by adding the taints [node-role.kubernetes.io/master:NoSchedule]
		[bootstrap-token] Using token: t2hu0m.iwosu060ldiaezuj
		[bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
		[bootstraptoken] configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
		[bootstraptoken] configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
		[bootstraptoken] configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
		[bootstraptoken] creating the "cluster-info" ConfigMap in the "kube-public" namespace
		[addons] Applied essential addon: CoreDNS
		[addons] Applied essential addon: kube-proxy
		
		Your Kubernetes master has initialized successfully!
		
		To start using your cluster, you need to run the following as a regular user:
		
		  mkdir -p $HOME/.kube
		  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
		  sudo chown $(id -u):$(id -g) $HOME/.kube/config
		
		You should now deploy a pod network to the cluster.
		Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
		  https://kubernetes.io/docs/concepts/cluster-administration/addons/
		
		You can now join any number of machines by running the following on each node
		as root:
		
		  kubeadm join 192.168.97.4:6443 --token t2hu0m.iwosu060ldiaezuj --discovery-token-ca-cert-hash sha256:7c9c9419b645f772338246fba984adf033a06add4e8583549de05a3ad504cd89
		
		`
		Expect(GetKubeJoinFromOutput(case1)).To(Equal("kubeadm join 192.168.97.4:6443 --token t2hu0m.iwosu060ldiaezuj --discovery-token-ca-cert-hash sha256:7c9c9419b645f772338246fba984adf033a06add4e8583549de05a3ad504cd89"))

	})
})
