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


function wait_calico () {
    echo "sleep 20s"
    sleep 20s
    timeout=10
    while true; do
        if [ $timeout == 0 ]; then
            echo "timeout waiting for calico started"
            exit 1
        fi 
        IFS=$'\n'
        quit="yes"
        for item in $(kubectl get pod -n kube-system | grep calico)
        do
            echo $item | grep Running || quit=no
        done
        sleep 5s
        if [ $quit == "yes" ]; then
            echo "calico is ready"
            break
        fi
        timeout=$((timeout-1))
        echo "calico is not ready"
    done
}

function wait_flannel () {
    echo "sleep 20s"
    sleep 20s
    timeout=5
    while true; do
        if [ $timeout == 0 ]; then
            echo "timeout waiting for flannel started"
            exit 1
        fi 
        IFS=$'\n'
        quit="yes"
        for item in $(kubectl get pod -n kube-system | grep flannel)
        do
            echo $item | grep Running || quit=no
        done
        sleep 5s
        if [ $quit == "yes" ]; then
            echo "calico is ready"
            break
        fi
        timeout=$((timeout-1))
        echo "calico is not ready"
    done
}
#do kubectl stuff here
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

MASTERIP=`hostname -I | cut -f1 -d" "`
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
        sed -i -e "s/{{etcd-key}}/$etcdkey/g; s/{{etcd-cert}}/$etcdcert/g; s/{{etcd-ca}}/$etcdca/g; s/{{masterip}}/$MASTERIP/g" ${CNIPATH}/calico/calico-etcd.yaml
        echo "apply yaml"
        kubectl apply -f ${CNIPATH}/calico/calico-etcd.yaml
        wait_calico
    else
        printf "mode %s do not support or not recoginzed" $MODE
    fi
fi

if [ $CNI == "flannel" ];then
    echo "apply yaml"
    kubectl apply -f ${CNIPATH}/${CNI}/${CNI}.yaml
    wait_flannel
fi 


