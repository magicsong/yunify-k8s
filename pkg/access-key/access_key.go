package key

import (
	"fmt"

	"github.com/yunify/qingcloud-sdk-go/config"
	"github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog"
)

type QingCloudAccessKeyHelper struct {
	Zone          string
	AccessKeyPath string
	userID        string

	qingCloudService *service.QingCloudService
	qingCloudConfig  *config.Config
}

func NewQingCloudAccessKeyHelper(zone string, filepath string) *QingCloudAccessKeyHelper {
	return &QingCloudAccessKeyHelper{
		Zone:          zone,
		AccessKeyPath: filepath,
	}
}

func (q *QingCloudAccessKeyHelper) Init() error {
	qcConfig, _ := config.NewDefault()
	if q.AccessKeyPath == "" {
		err := qcConfig.LoadUserConfig()
		if err != nil {
			return err
		}
	} else {
		err := qcConfig.LoadConfigFromFilepath(q.AccessKeyPath)
		if err != nil {
			return err
		}
	}
	q.qingCloudConfig = qcConfig
	qcService, err := service.Init(qcConfig)
	if err != nil {
		return err
	}
	q.qingCloudService = qcService
	api, _ := q.qingCloudService.Accesskey(q.Zone)
	output, err := api.DescribeAccessKeys(&service.DescribeAccessKeysInput{
		AccessKeys: []*string{&q.qingCloudConfig.AccessKeyID},
	})
	if err != nil {
		klog.Errorf("Failed to get userID")
		return err
	}
	if len(output.AccessKeySet) == 0 {
		err = fmt.Errorf("AccessKey %s have not userid", q.qingCloudConfig.AccessKeyID)
		return err
	}
	q.userID = *output.AccessKeySet[0].Owner
	return nil
}

func (q *QingCloudAccessKeyHelper) GetUserID() string {
	return q.userID
}

func (q *QingCloudAccessKeyHelper) GetService() *service.QingCloudService {
	return q.qingCloudService
}
