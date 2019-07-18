module github.com/magicsong/yunify-k8s

go 1.12

require (
	github.com/go-logr/logr v0.1.0 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/yunify/qingcloud-sdk-go v2.0.0-alpha.35.0.20190710082549-9b4f4db80863+incompatible
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.3
)

replace github.com/yunify/qingcloud-sdk-go v2.0.0-alpha.35.0.20190710082549-9b4f4db80863+incompatible => github.com/magicsong/qingcloud-sdk-go v2.0.0-alpha.33.0.20190712132900-6305afc1ddb5+incompatible
