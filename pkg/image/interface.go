package image

type Interface interface {
	CreateImageBasedInstanceID(string, string) (string, error)
	DeleteImage(...string) error
}
