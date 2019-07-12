package tag

import (
	"fmt"

	"github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog"
)

type qingcloudTagService struct {
	userID     string
	tagService *service.TagService
}

func NewQingCloudTagService(tagService *service.TagService, userID string) Interface {
	return &qingcloudTagService{
		userID:     userID,
		tagService: tagService,
	}
}

var _ Interface = &qingcloudTagService{}

func (q *qingcloudTagService) CreateTag(name string) (string, error) {
	color := RandomColor()
	input := &service.CreateTagInput{
		TagName: &name,
		Color:   &color,
	}
	output, err := q.tagService.CreateTag(input)
	if err != nil {
		return "", err
	}
	if *output.RetCode != 0 {
		err := fmt.Errorf("Error in creating instances tag, err: %s", *output.Message)
		return "", err
	}
	return *output.TagID, nil
}

func (q *qingcloudTagService) DeleteTag(id string) error {
	input := &service.DeleteTagsInput{
		Tags: []*string{&id},
	}
	output, err := q.tagService.DeleteTags(input)
	if err != nil {
		return err
	}
	if *output.RetCode != 0 {
		err := fmt.Errorf("Error in deleting  tag, err: %s", *output.Message)
		return err
	}
	return nil
}

func (q *qingcloudTagService) GetTagClusterByName(name string) (*TagCluster, error) {
	input := &service.DescribeTagsInput{
		SearchWord: &name,
		Verbose:    service.Int(1),
	}
	output, err := q.tagService.DescribeTags(input)
	if err != nil {
		klog.Error("Failed to initialize go sdk")
		return nil, err
	}
	if *output.RetCode != 0 {
		err := fmt.Errorf("Error in getting tag, err: %s", *output.Message)
		return nil, err
	}
	for _, tag := range output.TagSet {
		if *tag.Owner == q.userID && *tag.TagName == name {
			tagCluster := &TagCluster{
				TagID:     *tag.TagID,
				Instances: make([]string, 0),
			}
			for _, tagPair := range tag.ResourceTagPairs {
				if *tagPair.ResourceType == "instance" {
					tagCluster.Instances = append(tagCluster.Instances, *tagPair.ResourceID)
				}
			}
			return tagCluster, nil
		}
	}
	return nil, nil
}

func (q *qingcloudTagService) TagInstances(tagid string, instances []string) error {
	resourcePair := make([]*service.ResourceTagPair, len(instances))
	for index := 0; index < len(instances); index++ {
		resourcePair[index] = &service.ResourceTagPair{
			ResourceID:   &instances[index],
			ResourceType: service.String("instance"),
			TagID:        &tagid,
		}
	}
	input := &service.AttachTagsInput{
		ResourceTagPairs: resourcePair,
	}
	output, err := q.tagService.AttachTags(input)
	if err != nil {
		return err
	}
	if *output.RetCode != 0 {
		err := fmt.Errorf("Error in attaching tag, err: %s", *output.Message)
		return err
	}
	return nil
}
