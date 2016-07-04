package imageservice

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteToLocalStorageStoresImage(t *testing.T) {
	dir, err := ioutil.TempDir("", "imgStoreTest")
	if err != nil {
		t.Fatalf("error '%v' setting up image store for testing")
	}

	var localStore ImageStore
	localStore = LocalFileStore{Path: dir}

	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	blue := color.RGBA{0, 0, 255, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	fileName, err := localStore.SaveImage(img)

	if err != nil {
		t.Errorf("expected error from SaveImage to be '%v', was '%v'", nil, err)
	}

	file, err := os.Open(filepath.Join(dir, fileName))

	if err != nil {
		t.Errorf("expected error from opening saved file to be '%v', was '%v'", nil, err)
	}

	if _, _, err := image.Decode(file); err != nil {
		t.Errorf("expected error from decoding image to be '%v', was '%v'", nil, err)
	}
}
