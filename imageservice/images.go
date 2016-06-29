package imageservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"

	_ "image/jpeg"
	_ "image/png"

	proto "github.com/pivotal-sg/ochoku/imageservice/proto"
	"golang.org/x/net/context"
)

type ImageService struct{}

// Validation holds a field level valdiation error.  It also implements the error
// interface.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error stringifies the ValidationError, implements the error interface
func (err ValidationError) Error() string {
	return fmt.Sprintf("%s is %s", err.Field, err.Message)
}

func validateImageData(i *proto.ImageData) []error {
	errorSlice := make([]error, 0, 0)

	if i.Name == "" {
		errorSlice = append(errorSlice, ValidationError{Field: "name", Message: "missing"})
	}
	// An Image Without and Image...
	if len(i.Image) == 0 {
		errorSlice = append(errorSlice, ValidationError{Field: "image", Message: "missing"})
	}
	// If there is image data, make sure it is a valid image :)
	if len(i.Image) > 0 {
		buf := bytes.NewBuffer(i.Image)
		_, _, err := image.Decode(buf)
		if err == image.ErrFormat {
			errorSlice = append(errorSlice, ValidationError{Field: "image", Message: "Bad image format, jpeg and png only please."})
		}
	}

	if len(errorSlice) == 0 {
		return nil
	} else {
		return errorSlice
	}
}

func validateImageChoice(i *proto.ImageChoice) []error {
	// make sure we don't panic on pointer de-ref; and you shouldn't be validating
	// a nil anyway
	if i == nil {
		return []error{errors.New("Bad Data: nil")}
	}
	// name can't be blank.
	if i.Name == "" {
		return []error{ValidationError{Field: "name", Message: "missing"}}
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
		case ValidationError:
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

func (*ImageService) StoreImage(ctx context.Context, imgData *proto.ImageData, resp *proto.StatusResponse) error {
	buf := &bytes.Buffer{}
	buf.Write(imgData.Image)
	if errorSlice := validateImageData(imgData); errorSlice != nil {
		r, err := failedStatusResponse(errorSlice)
		*resp = *r
		return err
	}
	*resp = proto.StatusResponse{
		Message: "Image Saved!",
		Success: true,
	}
	return nil
}

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

func (*ImageService) ImagesFor(ctx context.Context, itemName *proto.ItemName, imageList *proto.ImageList) error {
	return nil
}
