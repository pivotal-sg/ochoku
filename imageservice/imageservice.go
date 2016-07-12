package imageservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"reflect"

	_ "image/jpeg"
	_ "image/png"

	"github.com/micro/go-platform/config"
	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
	"github.com/pivotal-sg/ochoku/validation"
	"golang.org/x/net/context"
)

type ImageService struct {
	Config config.Config

	DataStore []proto.ImageList
	FileStore ImageStore // TODO change this name.  it sucks
}

var imageDataValidations validation.Validations
var imageChoiceValidations validation.Validations

// Set up validations for the various types
func init() {
	imageDataValidations = validation.Validations{
		{
			V:         validation.ValidateByteSliceNotEmpty,
			FieldName: "Image",
		},
		{
			V:         validation.ValidateStringNotBlank,
			FieldName: "Name",
		},
		{
			V:         validateImageFormat,
			FieldName: "Image",
		},
	}

	imageChoiceValidations = validation.Validations{
		{
			V:         validation.ValidateStringNotBlank,
			FieldName: "Name",
		},
	}
}

// validateImageFormat ensure that the field dented by fieldName is
// one of the supported image formats.  It will return an error
// (ValidationError if its validation related) or another error
// for other things.
// It will not error on an empty byte array; you need to check that
// separatly.
func validateImageFormat(i interface{}, fieldName string) error {
	vb := reflect.ValueOf(i)
	if vb == reflect.Zero(reflect.TypeOf(i)) {
		return nil
	}

	if vb.FieldByName(fieldName).Type() != reflect.TypeOf([]byte{}) {
		return errors.New(fieldName + " is Not a []byte")
	}
	imgData := vb.FieldByName(fieldName).Bytes()

	// Ignore empty data sets: check that with `validation.ValidateByteSliceNotEmpty`
	if len(imgData) == 0 {
		return nil
	}

	buf := bytes.NewBuffer(imgData)
	_, _, err := image.Decode(buf)

	if err == image.ErrFormat {
		return validation.ValidationError{
			Field:   fieldName,
			Message: "Bad image format, jpeg and png only please.",
		}
		return err
	}
	return nil
}

func validateImageChoice(i *proto.ImageChoice) []error {
	// make sure we don't panic on pointer de-ref; and you shouldn't be validating
	// a nil anyway
	if i == nil {
		return []error{errors.New("Bad Data: nil")}
	}
	// name can't be blank.
	if i.Name == "" {
		return []error{validation.ValidationError{Field: "name", Message: "missing"}}
	}
	return nil
}

// failedStatusResponse converts a list of errors into a correct response.
// if they are expected errors (ValidationError), then it will return a StatusResponse,
// If they are unexpected, then it will return the error
func failedStatusResponse(errorSlice []error) (*proto.StatusResponse, error) {
	messageMap := make(map[string]string)
	for _, e := range errorSlice {
		switch err := e.(type) {
		case validation.ValidationError:
			messageMap[err.Field] = err.Message
		default:
			return nil, err
		}
	}
	msg, err := json.Marshal(messageMap)
	if err != nil {
		return nil, err
	}

	return &proto.StatusResponse{
		Message: string(msg),
		Success: false,
	}, nil
}

// StoreImage for the named item.  Should append to the currently stored images
// removal is handled by RemoveImage.  It will default to using the first image
// as the Cover image (which is a fancy way of saying, it lets you keep a pointer
// to which image you want to be the cover image)
func (is *ImageService) StoreImage(ctx context.Context, imgData *proto.ImageData, resp *proto.StatusResponse) error {
	buf := &bytes.Buffer{}
	buf.Write(imgData.Image)
	if errorSlice := imageDataValidations.Validate(*imgData); errorSlice != nil {
		r, err := failedStatusResponse(errorSlice)
		*resp = *r
		return err
	}

	img := &proto.Image{
		Caption: imgData.Caption,
		Uri:     "/1.jpeg",
	}

	imageList, index := is.getByName(imgData.Name)
	if index < 0 {
		imageList = &proto.ImageList{
			Name: imgData.Name,
			Images: []*proto.Image{
				img,
			},
		}
		is.DataStore = append(is.DataStore, *imageList)
	} else {
		imageList.Images = append(imageList.Images, img)
	}

	*resp = proto.StatusResponse{
		Message: "Image Saved!",
		Success: true,
	}
	return nil
}

// getByName the value and index of an image.  returns nil, -1 if not there.
func (is *ImageService) getByName(name string) (*proto.ImageList, int) {
	for i, imageList := range is.DataStore {
		if imageList.Name == name {
			return &imageList, i
		}
	}
	return nil, -1

}

// ChooseCover image for the named item passed in.  By default this will be the
// first image stored; because that is index 0.  This lets you change it.
func (*ImageService) ChooseCover(ctx context.Context, imgChoice *proto.ImageChoice, resp *proto.StatusResponse) error {
	errorSlice := validateImageChoice(imgChoice)
	if errorSlice != nil {
		r, err := failedStatusResponse(errorSlice)
		*resp = *r
		return err
	}
	resp.Message = "New Cover Selected"
	resp.Success = true
	return nil
}

func (*ImageService) RemoveImage(ctx context.Context, imgChoice *proto.ImageChoice, resp *proto.StatusResponse) error {
	return nil
}

func (is *ImageService) ImagesFor(ctx context.Context, itemName *proto.ItemName, imageList *proto.ImageList) error {
	for _, il := range is.DataStore {
		if il.Name == itemName.Name {
			*imageList = il
			return nil
		}
	}
	return nil
}
