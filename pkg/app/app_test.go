package app

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	It("Should be able to extract kubeadm join", func() {
		case1 := `
		You should now deploy a pod network to the cluster.
		Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
		  https://kubernetes.io/docs/concepts/cluster-administration/addons/
		
		You can now join any number of machines by running the following on each node
		as root:
		
		  kubeadm join 192.168.97.4:6443 --token t2hu0m.iwosu060ldiaezuj --discovery-token-ca-cert-hash sha256:7c9c9419b645f772338246fba984adf033a06add4e8583549de05a3ad504cd89
		
		`
		Expect(GetKubeJoinFromOutput(case1)).To(Equal("kubeadm join 192.168.97.4:6443 --token t2hu0m.iwosu060ldiaezuj --discovery-token-ca-cert-hash sha256:7c9c9419b645f772338246fba984adf033a06add4e8583549de05a3ad504cd89"))
		case1 = `You should now deploy a pod network to the cluster.
		Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
		  https://kubernetes.io/docs/concepts/cluster-administration/addons/
		
		Then you can join any number of worker nodes by running the following on each as root:
		
		kubeadm join 192.168.97.2:6443 --token ifqc4s.w1kemvf5d66v0qw1 \\
			--discovery-token-ca-cert-hash sha256:912f6349636027c61d5d98dbfef2393106119e47093efa721afe9522f963df32`
		Expect(GetKubeJoinFromOutput(case1)).To(Equal("kubeadm join 192.168.97.2:6443 --token ifqc4s.w1kemvf5d66v0qw1 --discovery-token-ca-cert-hash sha256:912f6349636027c61d5d98dbfef2393106119e47093efa721afe9522f963df32"))
	})
})
