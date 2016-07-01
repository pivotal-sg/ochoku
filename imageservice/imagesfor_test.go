package imageservice_test

import (
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

func TestGetOneImage(t *testing.T) {
	imageServiceObject := &imageservice.ImageService{
		Store: []proto.ImageList{
			{
				Name: "Choco A",
				Images: []*proto.Image{
					{
						Caption: "caption",
						Uri:     "/choco-a.jpeg",
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
						Uri:     "example.com/choco-a.jpg",
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
