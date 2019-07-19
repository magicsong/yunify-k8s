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
docker pull k8s.gcr.io/kube-proxy:v1.15.0
docker pull k8s.gcr.io/pause:3.1
docker pull k8s.gcr.io/coredns:1.3.1
docker pull quay.io/coreos/flannel:v0.11.0-amd64
docker pull calico/cni:v3.8.0
docker pull calico/pod2daemon-flexvol:v3.8.0
docker pull calico/node:v3.8.0
docker pull calico/kube-controllers:v3.8.0