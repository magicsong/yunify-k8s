package image

import (
	"fmt"
	"time"

	"github.com/magicsong/yunify-k8s/pkg/instance"
	"github.com/yunify/qingcloud-sdk-go/client"
	"github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog"
)

const DefaultCreateImageWait = time.Minute

type qingCloudImageService struct {
	instance.Interface
	imageService *service.ImageService
	jobService   *service.JobService
	userid       string
}

func (q *qingCloudImageService) CreateImageBasedInstanceID(instanceid string, imageName string) (string, error) {
	//stop instance
	err := q.StopInstances(instanceid)
	if err != nil {
		klog.Errorf("Failed to stop instance %s", instanceid)
		return "", err
	}
	input := &service.CaptureInstanceInput{
		ImageName: &imageName,
		Instance:  &instanceid,
	}
	output, err := q.imageService.CaptureInstance(input)
	if err != nil {
		klog.Error("error in capture instances, pls try again")
		return "", err
	}
	if *output.RetCode != 0 {
		err = fmt.Errorf("err: %s", *output.Message)
		klog.Error("error in capture  instances")
		return "", err
	}
	time.Sleep(time.Second * 30)
	klog.Info("Waiting for building image done")
	err = client.WaitJob(q.jobService, *output.JobID, DefaultCreateImageWait, time.Second*5)
	if err != nil {
		return "", err
	}
	return *output.ImageID, nil
}

func (q *qingCloudImageService) DeleteImage(ids ...string) error {
	input := &service.DeleteImagesInput{
		Images: service.StringSlice(ids),
	}
	output, err := q.imageService.DeleteImages(input)
	if err != nil {
		klog.Error("error in deleting images, pls try again")
		return err
	}
	if *output.RetCode != 0 {
		err = fmt.Errorf("err: %s", *output.Message)
		klog.Error(err, "error in delete images")
		return err
	}
	klog.Info("Waiting for image deletition done")
	err = client.WaitJob(q.jobService, *output.JobID, DefaultCreateImageWait, time.Second*5)
	if err != nil {
		return err
	}
	return nil
}

func NewQingCloudImageService(inst *service.InstanceService, job *service.JobService, image *service.ImageService, userid string) Interface {
	return &qingCloudImageService{
		jobService:   job,
		imageService: image,
		Interface:    instance.NewQingCloudInstanceService(inst, job),
	}
}
