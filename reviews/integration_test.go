package reviews_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/consul/command"
	"github.com/micro/go-platform/auth"
	"github.com/mitchellh/cli"
	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
)

type AuthTest struct {
	Authed bool
}

type tokenKey struct{}

func (a *AuthTest) Authorized(ctx context.Context, req auth.Request) (*auth.Token, error) {
	return a.Token()
}

func (a *AuthTest) Token() (*auth.Token, error) {
	if a.Authed {
		return &auth.Token{}, nil
	}
	return nil, auth.ErrInvalidToken
}

func (a *AuthTest) Introspect(ctx context.Context) (*auth.Token, error) {
	return a.Token()
}

func (a *AuthTest) Revoke(t *auth.Token) error {
	return nil
}

func (a *AuthTest) FromContext(ctx context.Context) (*auth.Token, bool) {
	return &auth.Token{}, a.Authed
}

func (a *AuthTest) NewContext(ctx context.Context, t *auth.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, t)
}

func (a *AuthTest) FromHeader(map[string]string) (*auth.Token, bool) {
	panic("not implemented")
}

func (a *AuthTest) NewHeader(map[string]string, *auth.Token) map[string]string {
	panic("not implemented")
}

func (a *AuthTest) String() string {
	return fmt.Sprintf("%v", a)
}

type integrationWrapper struct {
	a auth.Auth
}

func (i integrationWrapper) wrap(tests func(*testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		os.Remove("reviews_test.db")
		ui := &cli.MockUi{}
		consulReload := &command.ReloadCommand{Ui: ui}
		consulReload.Run([]string{})
		time.Sleep(100 * time.Millisecond)

		go func() {
			service := reviews.NewService("reviews_test.db", i.a)
			service.Run()
		}()

		time.Sleep(100 * time.Millisecond)
		tests(t)
	}
}

func TestIntegration(t *testing.T) {
	a := &AuthTest{true}
	wrapper := integrationWrapper{a: a}
	wrapper.wrap(func(t *testing.T) {
		client := reviews.NewClient()
		t.Run("storage=file", func(t *testing.T) {
			reviewRequest := &proto.ReviewRequest{
				Reviewer: "James",
				Name:     "Hershy's Dark Cardboard",
				Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
				Rating:   -5,
			}

			_, err := client.Review(context.TODO(), reviewRequest)

			if err != nil {
				t.Errorf("Expected error to be nil, was  '%v'", err)
			}

		})

		t.Run("storage=file2", func(t *testing.T) {
			expected := &proto.ReviewList{
				Count: 1,
				Reviews: []*proto.ReviewDetails{&proto.ReviewDetails{
					Reviewer: "James",
					Name:     "Hershy's Dark Cardboard",
					Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
					Rating:   -5,
				}}}
			allReviews, err := client.AllReviews(context.TODO(), &proto.Empty{})

			if err != nil {
				t.Errorf("Expected error to be nil, was  '%v'", err)
			}
			if !reflect.DeepEqual(expected, allReviews) {
				t.Errorf("Expected allReviews to be '%v', was  '%v'", expected, allReviews)
			}

		})
		t.Run("auth=fail", func(t *testing.T) {
			a.Authed = false
			reviewRequest := &proto.ReviewRequest{
				Reviewer: "James",
				Name:     "Hershy's Dark Cardboard",
				Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
				Rating:   -5,
			}

			_, err := client.Review(wrapper.a.NewContext(context.TODO(), &auth.Token{}), reviewRequest)

			if auth.ErrInvalidToken.Error() != err.Error() {
				t.Errorf("Expected error '%v', was  '%v'", auth.ErrInvalidToken, err)
			}
		})
	})(t)
}
