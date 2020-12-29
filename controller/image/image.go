package image

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"
	"path"
	"time"

	"image/jpeg"
	"image/png"

	"github.com/brxie/ebazarek-backend/config"
	"github.com/nfnt/resize"
)

type uploadResult struct {
	Path            string
	ThumbnailPath   string
	ThumbnailWidth  int
	ThumbnailHeight int
	Width           int
	Height          int
}

var uploadDir = config.Viper.GetString("UPLOAD_DIR")
var maxSize = uint(210)

func GetImage(fileName string) (*os.File, error) {
	file, err := os.Open(path.Join(uploadDir, fileName))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func UploadImage(fileName string, body []byte) (*uploadResult, error) {
	var img, thumb image.Image
	fileName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName)
	contentType := http.DetectContentType(body)

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// We don't want to exceed any dimension of the image by the specified value.
	// Depend on image is vertically or horizontally oriented, we should scale out
	// biggest side. Smallest size will be adjusted to the aspect ratio.
	if img.Bounds().Dx() > img.Bounds().Dy() {
		thumb = resize.Resize(maxSize, 0, img, resize.Lanczos3)
	} else {
		thumb = resize.Resize(0, maxSize, img, resize.Lanczos3)
	}

	result := &uploadResult{
		ThumbnailWidth:  thumb.Bounds().Dx(),
		ThumbnailHeight: thumb.Bounds().Dy(),
		Width:           img.Bounds().Dx(),
		Height:          img.Bounds().Dy(),
	}

	result.Path, err = writeImage(fileName, contentType, &img, false)
	if err != nil {
		return nil, err
	}

	result.ThumbnailPath, err = writeImage(fileName, contentType, &thumb, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func writeImage(originFileName, contentType string, image *image.Image, thumbnail bool) (string, error) {
	var err error
	if thumbnail {
		originFileName = "thumbnail_" + originFileName
	}
	path := path.Join(uploadDir, originFileName)
	fi, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	if contentType == "image/jpeg" {
		err = jpeg.Encode(fi, *image, nil)
	} else if contentType == "image/png" {
		err = png.Encode(fi, *image)
	} else {
		return "", fmt.Errorf(fmt.Sprintf("The '%s' content type is not supported", contentType))
	}

	return originFileName, err
}
