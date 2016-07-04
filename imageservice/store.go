package imageservice

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"

	_ "image/jpeg"
	_ "image/png"
)

type ImageStore interface {
	SaveImage(img image.Image) (string, error)
}

type LocalFileStore struct {
	Path string
}

func (ls LocalFileStore) SaveImage(img image.Image) (string, error) {
	filename := uuid.NewV4().String() + ".jpg"
	fd, err := os.Create(filepath.Join(ls.Path, filename))
	if err != nil {
		return "", err
	}
	if err := jpeg.Encode(fd, img, nil); err != nil {
		return filename, err
	}
	return filename, nil
}
