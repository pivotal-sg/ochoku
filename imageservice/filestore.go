package imageservice

import (
	"image"
	"os"
	"path/filepath"

	"github.com/micro/go-platform/config"
	"github.com/satori/go.uuid"
)

type LocalFileStore struct {
	// Path is the directory in which to store the image.
	Path   string
	Config config.Config
}

// SaveImage persists the image to the implemented storage backend,
// and returns the URI to said image.  For example, this version
// returns the path to the url, which can be combined with a hostname
// if needed, or just used as a relative import path.  This is mostly
// inteded to be used for testing.
// This version will store all images as a jpeg with default compression.
func (ls LocalFileStore) SaveImage(img image.Image) (string, error) {
	filename := uuid.NewV4().String() + ".jpg"
	fd, err := os.Create(filepath.Join(ls.Path, filename))
	if err != nil {
		return "", err
	}
	err = encodeToJpeg(img, fd)
	return filename, err
}
