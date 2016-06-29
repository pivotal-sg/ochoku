package imageservice_test

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"

	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
)

type ChooseCoverTestData struct {
	input    *proto.ImageChoice
	expected *proto.StatusResponse
}

func TestValidInputs(t *testing.T) {
	testData := []ChooseCoverTestData{
		{
			input: &proto.ImageChoice{
				Name:  "Choco A",
				Index: 0,
			},
			expected: &proto.StatusResponse{
				Message: "New Cover Selected",
				Success: true,
			},
		},
	}

	for _, testValue := range testData {
		var resp *proto.StatusResponse = &proto.StatusResponse{}
		if err := imageServiceObject.ChooseCover(context.TODO(), testValue.input, resp); err != nil {
			t.Errorf("for '%v'\nexpected error to be nil; was '%v'", testValue.input, err)
		}

		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("for '%v'\n\nexpected : '%v'\n\ngot      : '%v'", testValue.input, testValue.expected, resp)
		}
	}
}

func TestInvalidInputs(t *testing.T) {
	testData := []ChooseCoverTestData{
		{
			input: &proto.ImageChoice{
				Name:  "",
				Index: 0,
			},
			expected: &proto.StatusResponse{
				Message: `{"name":"missing"}`,
				Success: false,
			},
		},
	}

	for _, testValue := range testData {
		var resp *proto.StatusResponse = &proto.StatusResponse{}
		if err := imageServiceObject.ChooseCover(context.TODO(), testValue.input, resp); err != nil {
			t.Errorf("for '%v'\nexpected error to be nil; was '%v'", testValue.input, err)
		}

		if !reflect.DeepEqual(testValue.expected, resp) {
			t.Errorf("for '%v'\n\nexpected : '%v'\n\ngot      : '%v'", testValue.input, testValue.expected, resp)
		}
	}
}
