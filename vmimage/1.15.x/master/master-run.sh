#!/bin/bash

swapoff -a
sed '/exit/i\swapoff -a' -i /etc/rc.local
sysctl net.bridge.bridge-nf-call-iptables=1

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

## Download GPG key.
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
## Add docker apt repository.
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

apt-get update && apt-get install -y apt-transport-https ca-certificates curl software-properties-common

apt-get install -y kubelet kubeadm kubectl docker.io jq
apt-mark hold kubelet kubeadm kubectl

echo "source <(kubectl completion bash)" >> ~/.bashrc
## Install docker ce.
#apt-get install docker-ce=18.06.2~ce~3-0~ubuntu

# Setup daemon.
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

mkdir -p /etc/systemd/system/docker.service.d

# Restart docker.
systemctl daemon-reload
systemctl enable docker
systemctl restart docker


##pull image
kubeadm config images list
kubeadm config images pull

##pull CNI image
mkdir -p CNI/flannel
wget https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml -O CNI/flannel/kube-flannel.yml
docker pull quay.io/coreos/flannel:v0.11.0-amd64

mkdir -p CNI/calico
wget https://docs.projectcalico.org/v3.8/manifests/calico.yaml -O CNI/calico/calico.yaml
wget https://raw.githubusercontent.com/magicsong/yunify-k8s/master/vmimage/1.15.x/master/calico-etcd.yaml -O CNI/calico/calico-etcd.yaml
docker pull calico/cni:v3.8.1
docker pull calico/pod2daemon-flexvol:v3.8.1
docker pull calico/node:v3.8.1
docker pull calico/kube-controllers:v3.8.1

#install calicoctl
install_calicoctl=0
command -v calicoctl >/dev/null 2>&1 ||  install_calicoctl=1
if [ $install_calicoctl == 1 ]; then
  curl -O -L  https://github.com/projectcalico/calicoctl/releases/download/v3.8.1/calicoctl
  chmod +x calicoctl
  mv calicoctl /usr/local/bin/
fi

mkdir -p /etc/calico
cat <<EOF > /etc/calico/calicoctl.cfg
apiVersion: projectcalico.org/v3
kind: CalicoAPIConfig
metadata:
spec:
  datastoreType: "etcdv3"
  etcdEndpoints: "https://localhost:2379"
  etcdKeyFile: /etc/kubernetes/pki/apiserver-etcd-client.key
  etcdCertFile: /etc/kubernetes/pki/apiserver-etcd-client.crt
  etcdCACertFile: /etc/kubernetes/pki/etcd/ca.crt
EOF