#!/bin/bash
POD_CIDR="192.168.0.0/16"
CNI="calico"
CNIPATH=/root/CNI
MODE="k8s"

set -e
# parse args
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    --pod-cidr)
    POD_CIDR=$2
    shift
    shift # past argument
    ;;
    -n|--CNI)
    CNI=$2
    shift # past argument
    shift # past value
    ;;
    -m|--mode)
    MODE=$2
    shift # past argument
    shift # past value
    ;;
    -t|--tag)
    tag="$2"
    shift # past argument
    shift # past value
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done

#do kubectl stuff here
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

echo "CNI=${CNI}, pod-cidr=${POD_CIDR}"

if [ $CNI == "calico" ];then
    if [ $MODE == "k8s" ]; then
        if [  $POD_CIDR != "192.168.0.0/16" ]; then
            sed -i -e "s?192.168.0.0/16?$POD_CIDR?g" ${CNIPATH}/calico/calico.yaml
        fi
        echo "apply yaml"
        kubectl apply -f ${CNIPATH}/calico/calico.yaml           
    elif [ $MODE == "etcd" ]; then
        if [  $POD_CIDR != "192.168.0.0/16" ]; then
            sed -i -e "s?192.168.0.0/16?$POD_CIDR?g" ${CNIPATH}/calico/calico-etcd.yaml
        fi
        # set etcd ca
        etcdkey=`cat /etc/kubernetes/pki/apiserver-etcd-client.key | base64 -w 0`
        etcdcert=$(cat /etc/kubernetes/pki/apiserver-etcd-client.crt | base64 -w 0)
        etcdca=$(cat /etc/kubernetes/pki/etcd/ca.crt | base64 -w 0)
        sed -i -e "s/{{etcd-key}}/$etcdkey/g" ${CNIPATH}/calico/calico-etcd.yaml
        sed -i -e "s/{{etcd-cert}}/$etcdcert/g" ${CNIPATH}/calico/calico-etcd.yaml
        sed -i -e "s/{{etcd-ca}}/$etcdca/g" ${CNIPATH}/calico/calico-etcd.yaml

        echo "apply yaml"
        kubectl apply -f ${CNIPATH}/calico/calico-etcd.yaml
    else
        printf "mode %s do not support or not recoginzed" $MODE
    fi
fi

if [ $CNI == "flannel" ];then
    echo "apply yaml"
    kubectl apply -f ${CNIPATH}/${CNI}/${CNI}.yaml
fi 


