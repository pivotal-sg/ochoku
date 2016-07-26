package userservice

import (
  "time"

  proto "github.com/pivotal-sg/ochoku/userservice/proto"
  "github.com/pivotal-sg/ochoku/userservice/storage"
  "golang.org/x/crypto/bcrypt"
  "golang.org/x/net/context"
)

type UserService struct {
  Store storage.Storer
}

func (us *UserService) Register(ctx context.Context, regData *proto.RegistrationData, userDetails *proto.UserDetails) error {
  hash, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
  if err != nil {
    return err
  }

  userDetails.Username = regData.Username
  userDetails.HashedPassword = string(hash)
  userDetails.Joined = time.Now().Unix()
  userDetails.Name = regData.Name

  return us.Store.Insert(*userDetails)
}

func (us *UserService) PasswordLogin(ctx context.Context, loginDetails *proto.LoginDetails, loginStatus *proto.LoginStatus) error {
  user, err := us.Store.Get(loginDetails.Username)
  if err != nil {
    return err
  }

  if hashErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginDetails.Password)); hashErr != nil {
    loginStatus.Ok = false
    return nil
  }
  loginStatus.Ok = true
  return nil
}
