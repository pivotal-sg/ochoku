package userservice_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pivotal-sg/ochoku/userservice"
	proto "github.com/pivotal-sg/ochoku/userservice/proto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type MockStore struct {
	users      map[string]proto.UserDetails
	throwError bool
}

func (ms *MockStore) Get(username string) (proto.UserDetails, error) {
	if ms.throwError {
		return proto.UserDetails{}, errors.New("I fail it")
	}
	return ms.users[username], nil
}

func (ms *MockStore) Insert(usr proto.UserDetails) error {
	ms.users[usr.Username] = usr
	return nil
}

func (ms *MockStore) reset() {
	ms.users = make(map[string]proto.UserDetails)
	ms.throwError = false
}

var mockStore *MockStore = &MockStore{users: make(map[string]proto.UserDetails)}
var userServiceObject *userservice.UserService = &userservice.UserService{
	Store: mockStore,
}

func createUser(ms *MockStore) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	ms.Insert(proto.UserDetails{
		Username:       "username",
		HashedPassword: string(hash),
		Joined:         12345,
		Name:           "Joan Smith"})
}

func TestUserRegistrationReturnsDetails(t *testing.T) {
	mockStore.reset()
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

	if response.Joined == 0 || response.Joined > time.Now().Unix() {
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
	mockStore.reset()
	createUser(mockStore)

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

func TestErrorInGetOnLoginReturnsError(t *testing.T) {
	mockStore.reset()
	mockStore.throwError = true

	var response *proto.LoginStatus = &proto.LoginStatus{}
	ctx := context.TODO()

	loginRequest := proto.LoginDetails{
		Username: "username",
		Password: "password",
	}

	err := userServiceObject.PasswordLogin(ctx, &loginRequest, response)

	if err == nil {
		t.Errorf("expected the response to have an error, it was %v", err)
	}
}
