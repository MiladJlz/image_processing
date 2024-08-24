package types

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"image_processing/errors"
)

type Status string

const (
	bcryptCost     = 12
	minUsernameLen = 6
	maxUsernameLen = 15
	minPasswordLen = 7
)

type User struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encryptedPassword"`
}
type CreateUserParams struct {
	Username string
	Password string
}

func NewUserFromParams(params CreateUserParams) (*User, *errors.Error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, errors.ErrServer(err.Error())
	}
	return &User{
		Username:          params.Username,
		EncryptedPassword: string(encpw),
	}, nil
}
func (params CreateUserParams) Validate() *errors.Error {

	if len(params.Username) < minUsernameLen || len(params.Username) > maxUsernameLen {
		return errors.ErrBadRequest(fmt.Sprintf("username length should be between  %d and %d characters", minUsernameLen, maxUsernameLen))
	}

	if len(params.Password) < minPasswordLen {
		return errors.ErrBadRequest(fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}
	return nil
}
