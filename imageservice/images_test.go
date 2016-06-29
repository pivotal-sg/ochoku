package images_test

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/pivotal-sg/ochoku/images"
	proto "github.com/pivotal-sg/ochoku/images/proto"
)

var imageServiceObject images.ImageService = images.ImageService{}

type encoder func(io.Writer, image.Image) error

func testImage(encoder func(io.Writer, image.Image) error) []byte {
	m := image.NewRGBA(image.Rect(0, 0, 50, 50))
	blue := color.RGBA{0, 0, 255, 255}

	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	buf := &bytes.Buffer{}
	encoder(buf, m)
	return buf.Bytes()
}

func TestUploadImageReturnsSuccessForValidFormats(t *testing.T) {
	var testData = []struct {
		input    *proto.ImageData
		expected *proto.StatusResponse
	}{
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as jpeg",
				Image:   testImage(func(w io.Writer, i image.Image) error { return jpeg.Encode(w, i, nil) }),
			},
			expected: &proto.StatusResponse{
				Message: "Image Saved!",
				Success: true,
			},
		},
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as png",
				Image:   testImage(png.Encode),
			},
			expected: &proto.StatusResponse{
				Message: "Image Saved!",
				Success: true,
			},
		},
	}

	for _, testValue := range testData {
		var resp = &proto.StatusResponse{}
		if err := imageServiceObject.StoreImage(context.TODO(), testValue.input, resp); err != nil {
			t.Fatalf(err.Error())
		}
		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("For input '%v': expected '%v', got '%v'", testValue.input, testValue.expected, resp)
		}
	}

}

func TestUploadImageReturnsFailForInValidImageFormats(t *testing.T) {
	var testData = []struct {
		input    *proto.ImageData
		expected *proto.StatusResponse
	}{
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as jpeg",
				Image:   []byte{0, 0, 0},
			},
			expected: &proto.StatusResponse{
				Message: "Bad image format, jpeg and png only please.",
				Success: false,
			},
		},
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as png",
				Image:   []byte{0, 0, 0},
			},
			expected: &proto.StatusResponse{
				Message: "Bad image format, jpeg and png only please.",
				Success: false,
			},
		},
	}

	for _, testValue := range testData {
		var resp = &proto.StatusResponse{}
		if err := imageServiceObject.StoreImage(context.TODO(), testValue.input, resp); err != nil {
			t.Fatalf(err.Error())
		}
		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("For input '%v': expected '%v', got '%v'", testValue.input, testValue.expected, resp)
		}
	}

}
