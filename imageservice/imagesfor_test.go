package imageservice_test

import (
	"fmt"
	"image"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/pivotal-sg/ochoku/imageservice"
	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
)

type ImagesForTestData struct {
	input    *proto.ItemName
	expected *proto.ImageList
}

type MockImageStore struct {
	I int
}

func (is *MockImageStore) SaveImage(img image.Image) (filename string, err error) {
	filename = fmt.Sprintf("/%d.jpg", is.I)
	is.I++
	return
}

func TestGetOneImage(t *testing.T) {
	imageServiceObject := &imageservice.ImageService{
		FileStore: &MockImageStore{I: 1},
		DataStore: []proto.ImageList{
			{
				Name: "Choco A",
				Images: []*proto.Image{
					{
						Caption: "caption",
						Uri:     "/1.jpeg",
					},
				},
			},
		},
	}
	testData := []ImagesForTestData{
		{
			input: &proto.ItemName{
				Name: "Choco A",
			},
			expected: &proto.ImageList{
				Name:  "Choco A",
				Cover: 0,
				Images: []*proto.Image{
					{
						Uri:     "/1.jpeg",
						Caption: "caption",
					},
				},
			},
		},
	}

	for _, testValue := range testData {
		var resp *proto.ImageList = &proto.ImageList{}
		if err := imageServiceObject.ImagesFor(context.TODO(), testValue.input, resp); err != nil {
			t.Errorf("for '%v'\nexpected error to be nil; was '%v'", testValue.input, err)
		}

		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("for '%v'\n\nexpected : '%v'\n\ngot      : '%v'", testValue.input, testValue.expected, resp)
		}
	}
}
