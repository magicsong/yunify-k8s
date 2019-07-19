#!/bin/bash
POD_CIDR="192.168.0.0/16"
CNI="calico"
CNIPATH=/root/CNI

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

echo "CNI=${CNI}, pod-cidr=${POD_CIDR}"

if [ $CNI == "calico" -a $POD_CIDR != "192.168.0.0/16" ]; then
    sed -i -e "s?192.168.0.0/16?$POD_CIDR?g" ${CNIPATH}/calico/calico.yaml
fi

echo "apply yaml"
kubectl apply -f ${CNIPATH}/${CNI}/${CNI}.yaml