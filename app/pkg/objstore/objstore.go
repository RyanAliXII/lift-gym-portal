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
func (s *ObjectStorage) Upload(ctx context.Context, file multipart.File, folderName string,  filename string ) (string, error) {
	result, err := s.cld.Upload.Upload(ctx, file,  uploader.UploadParams{
		Folder: folderName,
		PublicID: filename,
	})
	return result.PublicID, err
}
func (s *ObjectStorage)Remove(filepath string) error {
	var  invalidate bool = true
	_, err := s.cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: filepath,
		ResourceType: "image",
		Invalidate: &invalidate ,
	})
	return err
}
type ObjectStorer interface {
	Upload(ctx context.Context, file multipart.File, folderName string,  filename string ) (string, error)
	Remove(filename string) error
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


