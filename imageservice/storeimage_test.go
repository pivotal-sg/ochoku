package imageservice_test

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/pivotal-sg/ochoku/imageservice"
	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
)

var imageServiceObject imageservice.ImageService = imageservice.ImageService{}

type encoder func(io.Writer, image.Image) error

func testImage(encoder func(io.Writer, image.Image) error) []byte {
	m := image.NewRGBA(image.Rect(0, 0, 50, 50))
	blue := color.RGBA{0, 0, 255, 255}

	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
	buf := &bytes.Buffer{}
	encoder(buf, m)
	return buf.Bytes()
}

type StoreImageTestData struct {
	input    *proto.ImageData
	expected *proto.StatusResponse
}

// checkStoreResponse iterates over a table test for StoreImage
func checkStoreResponse(t *testing.T, testData []StoreImageTestData) {
	dir, _ := ioutil.TempDir("", "imgserviceTest")
	serviceObject := imageservice.ImageService{
		DataStore: make([]proto.ImageList, 0, 0),
		FileStore: imageservice.LocalFileStore{Path: dir},
	}
	for _, testValue := range testData {
		var resp = &proto.StatusResponse{}
		if err := serviceObject.StoreImage(context.TODO(), testValue.input, resp); err != nil {
			t.Fatalf(err.Error())
		}
		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("For input '%v': expected '%v', got '%v'", testValue.input, testValue.expected, resp)
		}
	}
}

func TestUploadImageReturnsSuccessForValidFormats(t *testing.T) {
	var testData = []StoreImageTestData{
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

	checkStoreResponse(t, testData)

}

func TestUploadImageReturnsFailForInValidImageFormats(t *testing.T) {
	var testData = []StoreImageTestData{
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as jpeg",
				Image:   []byte{0, 0, 0},
			},
			expected: &proto.StatusResponse{
				Message: `{"Image":"Bad image format, jpeg and png only please."}`,
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
				Message: `{"Image":"Bad image format, jpeg and png only please."}`,
				Success: false,
			},
		},
	}

	checkStoreResponse(t, testData)
}

func TestCaptionIsOptional(t *testing.T) {
	var testData = []StoreImageTestData{
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "",
				Image:   testImage(png.Encode),
			},
			expected: &proto.StatusResponse{
				Message: "Image Saved!",
				Success: true,
			},
		},
	}

	checkStoreResponse(t, testData)
}

func TestUploadImageReturnsFailForInValidFieldValues(t *testing.T) {
	var testData = []StoreImageTestData{
		{
			input: &proto.ImageData{
				Name:    "",
				Caption: "A Picture of 'Choco A' as png",
				Image:   testImage(png.Encode),
			},
			expected: &proto.StatusResponse{
				Message: `{"Name":"missing"}`,
				Success: false,
			},
		},
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "A Picture of 'Choco A' as jpeg",
				Image:   []byte{},
			},
			expected: &proto.StatusResponse{
				Message: `{"Image":"missing"}`,
				Success: false,
			},
		},
		{
			input: &proto.ImageData{
				Name:    "",
				Caption: "A Picture of 'Choco A' as jpeg",
				Image:   []byte{},
			},
			expected: &proto.StatusResponse{
				Message: `{"Image":"missing","Name":"missing"}`,
				Success: false,
			},
		},
	}
	checkStoreResponse(t, testData)
}

func TestUploadSavesTheURIBasedOnImageStorersReturn(t *testing.T) {
	serviceObject := imageservice.ImageService{
		DataStore: make([]proto.ImageList, 0, 0),
		FileStore: &MockImageStore{I: 1},
	}

	var testData = []StoreImageTestData{
		{
			input: &proto.ImageData{
				Name:    "Choco A",
				Caption: "caption",
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
		if err := serviceObject.StoreImage(context.TODO(), testValue.input, resp); err != nil {
			t.Fatalf(err.Error())
		}
		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("For input '%v': expected '%v', got '%v'", testValue.input, testValue.expected, resp)
		}
	}

	input := &proto.ItemName{
		Name: "Choco A",
	}
	expected := &proto.ImageList{
		Name:  "Choco A",
		Cover: 0,
		Images: []*proto.Image{
			{
				Uri:     "/1.jpeg",
				Caption: "caption",
			},
		},
	}

	var resp *proto.ImageList = &proto.ImageList{}
	if err := serviceObject.ImagesFor(context.TODO(), input, resp); err != nil {
		t.Errorf("for '%v'\nexpected error to be nil; was '%v'", input, err)
	}

	if !reflect.DeepEqual(expected, resp) {
		t.Errorf("for '%v'\n\nexpected : '%v'\n\ngot      : '%v'", input, expected, resp)
	}
}
