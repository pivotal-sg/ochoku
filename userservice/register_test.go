package userservice_test

import (
	"testing"
	"time"

	"github.com/pivotal-sg/ochoku/userservice"
	proto "github.com/pivotal-sg/ochoku/userservice/proto"
	"golang.org/x/net/context"
)

var userServiceObject *userservice.UserService = &userservice.UserService{}

func createUser(us *userservice.UserService) {
}

func TestUserRegistrationReturnsDetails(t *testing.T) {
	var response *proto.UserDetails = &proto.UserDetails{}
	ctx := context.TODO()

	registerRequest := proto.RegistrationData{
		Username: "username",
		Password: "password",
		Name:     "Joan Smith",
	}

	err := userServiceObject.Register(ctx, &registerRequest, response)

	if err != nil {
		t.Errorf("expected the response to not have an error, it was %v", err)
	}

	if response.Joined == 0 || response.Joined >= time.Now().Unix() {
		t.Errorf("expected a valid Joined time, was %v", response.Joined)
	}

	if response.Username != "username" {
		t.Errorf("expected Username to be '%v', was '%v'", registerRequest.Username, response.Username)
	}

	if response.Name != "Joan Smith" {
		t.Errorf("expected Username to be '%v', was '%v'", registerRequest.Name, response.Name)
	}
}

func TestUserLoginWorks(t *testing.T) {
	createUser(userServiceObject)

	var response *proto.LoginStatus = &proto.LoginStatus{}
	ctx := context.TODO()

	loginRequest := proto.LoginDetails{
		Username: "username",
		Password: "password",
	}

	err := userServiceObject.PasswordLogin(ctx, &loginRequest, response)

	if err != nil {
		t.Errorf("expected the response to not have an error, it was %v", err)
	}

	if !response.Ok {
		t.Errorf("expected Login status ok to be true.  wasn't")
	}

	if response.Msg != "" {
		t.Errorf("expected Msg to be blank, was '%v'", response.Msg)
	}

}
