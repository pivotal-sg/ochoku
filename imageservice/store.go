package imageservice

import (
	"image"
	"io"

	"image/jpeg"
)

// ImageStore is an interface to a backend for storing images.
type ImageStore interface {
	SaveImage(img image.Image) (string, error)
}

// encodeToJpeg is a helper function to encode the image to a
// what ever writer is needed for the storage backend.
func encodeToJpeg(img image.Image, w io.Writer) error {
	if err := jpeg.Encode(w, img, nil); err != nil {
		return err
	}
	return nil
}
