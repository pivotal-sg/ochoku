package userservice

import (
	proto "github.com/pivotal-sg/ochoku/userservice/proto"
	"golang.org/x/net/context"
)

type UserService struct{}

func (us *UserService) Register(ctx context.Context, regData *proto.RegistrationData, userDetails *proto.UserDetails) error {

	return nil
}

func (us *UserService) PasswordLogin(ctx context.Context, loginDetails *proto.LoginDetails, loginStatus *proto.LoginStatus) error {
	return nil
}
