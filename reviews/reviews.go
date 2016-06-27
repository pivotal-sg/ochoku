package reviews

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/micro/go-micro/client"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

type ReviewService struct{}

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

// failedStatusResponse converts a list of errors into a correct response.
// if they are expected errors (ValidationError), then it will return a StatusResponse,
// If they are unexpected, then it will return the error
func failedStatusResponse(errors []error) (*proto.StatusResponse, error) {
	messageMap := make(map[string]string)
	for _, e := range errors {
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

// validateName is not blank
func validateName(request proto.ReviewRequest) error {
	if request.Name == "" {
		return ValidationError{Field: "name", Message: "missing"}
	}
	return nil
}

// validateReviewer is not blank
func validateReviewer(request proto.ReviewRequest) error {
	if request.Reviewer == "" {
		return ValidationError{Field: "reviewer", Message: "missing"}
	}
	return nil
}

// Review a product, it should have a Name and Reviewer.
// It will return a StatusResponse as long as we know how to deal with what was passed in (Eg, known
// invalid data), or an error if something else was wrong (like ?)
func (*ReviewService) Review(c context.Context, reviewRequest *proto.ReviewRequest, opts ...client.CallOption) (*proto.StatusResponse, error) {
	if reviewRequest == nil {
		return nil, errors.New("ReviewRequest was nil, must be valid reference")
	}
	errors := make([]error, 0)

	if err := validateName(*reviewRequest); err != nil {
		errors = append(errors, err)
	}
	if err := validateReviewer(*reviewRequest); err != nil {
		errors = append(errors, err)
	}

	if len(errors) != 0 {
		return failedStatusResponse(errors)
	}

	return &proto.StatusResponse{
		Message: "All Good!",
		Success: true,
	}, nil
}
