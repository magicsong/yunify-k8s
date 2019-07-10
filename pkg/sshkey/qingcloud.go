package sshkey

import (
	"fmt"

	"github.com/yunify/qingcloud-sdk-go/service"
)

type qingcloudSSHKey struct {
	keyPairService *service.KeyPairService
	userID         string
}

func NewQingCloudKeyPairService(keyPairService *service.KeyPairService, userid string) Interface {
	return &qingcloudSSHKey{
		keyPairService: keyPairService,
		userID:         userid,
	}
}

func (q *qingcloudSSHKey) GetKeyPairByName(name string) (string, error) {
	input := &service.DescribeKeyPairsInput{
		SearchWord: &name,
		Owner:      &q.userID,
	}
	output, err := q.keyPairService.DescribeKeyPairs(input)
	if err != nil {
		return "", err
	}
	if *output.RetCode != 0 {
		err = fmt.Errorf("Error in getting ssh keypair by name %s, err: %s", name, *output.Message)
		return "", err
	}
	for _, key := range output.KeyPairSet {
		if *key.KeyPairName == name {
			return *key.KeyPairID, nil
		}
	}
	return "", nil
}

func (q *qingcloudSSHKey) CreateSSHKey(name string, key string) (string, error) {
	input := &service.CreateKeyPairInput{
		Mode:        service.String("user"),
		PublicKey:   &key,
		KeyPairName: &name,
	}
	output, err := q.keyPairService.CreateKeyPair(input)
	if err != nil {
		return "", err
	}
	if *output.RetCode != 0 {
		err = fmt.Errorf("Error in creating ssh keypair,err: %s", *output.Message)
		return "", err
	}

	return *output.KeyPairID, nil
}

func (q *qingcloudSSHKey) DeleteSSHKey(id string) error {
	input := &service.DeleteKeyPairsInput{
		KeyPairs: []*string{&id},
	}
	output, err := q.keyPairService.DeleteKeyPairs(input)
	if err != nil {
		return err
	}
	if *output.RetCode != 0 {
		err = fmt.Errorf("Error in creating ssh keypair,err: %s", *output.Message)
		return err
	}
	return nil
}
