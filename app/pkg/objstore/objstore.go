package objstore

import (
	"os"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
)

type ObjectStorage struct {
	cld * cloudinary.Cloudinary
}

type ObjectStorer interface {
}

func (s *ObjectStorage) Upload() {

}

var objecStorage ObjectStorage;
var initErr error;
var once sync.Once;
func getObjectStorage() (ObjectStorer, error) {
	once.Do(func() {
		cloudName := os.Getenv("CLOUDINARY_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
		objecStorage.cld, initErr = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
		
	})
	return objecStorage, nil
}