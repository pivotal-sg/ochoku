package imageservice

import (
	"image"

	_ "image/jpeg"
	_ "image/png"
)

// ImageStore is an interface to a backend for storing images.
type ImageStore interface {
	SaveImage(img image.Image) (string, error)
}
