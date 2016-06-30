package reviews

import (
	"encoding/json"
	"errors"

	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"github.com/pivotal-sg/ochoku/reviews/storage"
	"github.com/pivotal-sg/ochoku/validation"
	"golang.org/x/net/context"
)

type ReviewService struct {
	Store storage.Storer
}

var reviewValidations validation.Validations

func init() {
	reviewValidations = validation.Validations{
		{
			V:         validation.ValidateStringNotBlank,
			FieldName: "Name",
		},
		{
			V:         validation.ValidateStringNotBlank,
			FieldName: "Reviewer",
		},
	}
}

// failedStatusResponse converts a list of errors into a correct response.
// if they are expected errors (ValidationError), then it will return a StatusResponse,
// If they are unexpected, then it will return the error
func failedStatusResponse(errors []error) (*proto.StatusResponse, error) {
	messageMap := make(map[string]string)
	for _, e := range errors {
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

// Review a product, it should have a Name and Reviewer.
// It will return a StatusResponse as long as we know how to deal with what was passed in (Eg, known
// invalid data), or an error if something else was wrong (like ?)
func (rs *ReviewService) Review(c context.Context, reviewRequest *proto.ReviewRequest, response *proto.StatusResponse) error {
	if reviewRequest == nil {
		return errors.New("ReviewRequest was nil, must be valid reference")
	}

	if rs.Store == nil {
		return errors.New("Storer not set in context or wrong type")
	}

	errors := reviewValidations.Validate(*reviewRequest)

	// If i have errors: fail it
	if len(errors) != 0 {
		res, err := failedStatusResponse(errors)
		*response = *res
		return err
	}

	reviewDetails := proto.ReviewDetails{
		Name:     reviewRequest.Name,
		Reviewer: reviewRequest.Reviewer,
		Review:   reviewRequest.Review,
		Rating:   reviewRequest.Rating,
	}

	rs.Store.Insert(reviewDetails)

	*response = proto.StatusResponse{
		Message: "All Good!",
		Success: true,
	}
	return nil
}

// AllReviews will return all of the reviews so far
func (rs *ReviewService) AllReviews(context context.Context, empty *proto.Empty, response *proto.ReviewList) error {
	if rs.Store == nil {
		return errors.New("Storer not set in context or wrong type")
	}
	allReviews, _ := rs.Store.List()
	*response = proto.ReviewList{Reviews: allReviews, Count: int32(len(allReviews))}

	return nil
}
