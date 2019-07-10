package sshkey

type Interface interface {
	CreateSSHKey(string, string) (string, error)
	DeleteSSHKey(string) error
	GetKeyPairByName(string) (string, error)
}
