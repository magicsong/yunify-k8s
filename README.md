# yunify_k8s

> 个人项目

由于青云目前正在全力打造企业级容器引擎[Kubesphere](https://github.com/kubesphere/kubesphere)，导致其AppCenter上的原生k8s已经年久失修，最新的版本是1.10，而1.10的已经不再k8s社区的维护生命周期内，容易出现一些BUG和安全问题。Kubesphere支持最新的版本，但是由于其本身过于庞大，安装时间漫长，不利于开发人员快速起一个测试环境，所以开发了这个小工具。

## 功能

1. 2分钟后内起一个三节点集群，并且配置好网络插件
2. 20s 删除一个集群
3. 集群镜像定制（WIP）
4. 重启不会删除机器，放心使用本地文件

## 使用准备条件
1. 准备青云AccessKey文件，参考[官方文档](https://docs.qingcloud.com/product/cli/#%E6%96%B0%E6%89%8B%E6%8C%87%E5%8D%97)，将配置文件放在适当的位置。这是创建机器的凭证。
2. 在青云平台上创建VPC，并且通过VPN连接到VPC中。因为新创的机器没有公网IP，所以需要用VPN通过内网ip的方式访问集群机器。配置VPN请参考[官方文档](https://docs.qingcloud.com/product/network/vpn)
3. 本地已有SSH公钥，在`$HOME/.ssh/id_rsa.pub`，目前只支持这么一种SSH

## 使用方式

1. 从release页面下载最新binary，将其放入`$Path`中
2. 创建集群，最小参数需要指定集群所在Vxnet
```bash
# 默认配置是一个4核4G的基础型master和两个4核4G的基础型node，网络插件为calico，k8s版本为1.13.1
qks create cluster testk8s -x=vxnet-xxx
# 完整的用法请使用`qks create -h`
```
3. 删除集群
```bash
qks delete cluster testk8s
```

## 目前支持的版本
+ 1.13.x
+ 1.15.0