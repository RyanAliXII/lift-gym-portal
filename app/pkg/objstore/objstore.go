package objstore

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type ObjectStorage struct {
	cld * cloudinary.Cloudinary
}
func (s *ObjectStorage) Upload(file multipart.File, folderName string,  filename string ) (error) {
	_, err := s.cld.Upload.Upload(context.Background(), file,  uploader.UploadParams{
		Folder: folderName,
		PublicID: filename,
	})
	return err
}

type ObjectStorer interface {
	Upload(file multipart.File, folderName string,  filename string ) (error)
}
var PublicURL string;
var objecStorage * ObjectStorage;
var initErr error;
var once sync.Once;
func GetObjectStorage() (ObjectStorer, error) {
	once.Do(func() {
		cloudName := os.Getenv("CLOUDINARY_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
		objecStorage = &ObjectStorage{}
		PublicURL = fmt.Sprintf("https://res.cloudinary.com/%s", cloudName)
		objecStorage.cld, initErr = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
		
	})
	return objecStorage, initErr
}


