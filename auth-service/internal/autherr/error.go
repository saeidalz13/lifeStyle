package autherr

import "errors"

var (
	ErrAuthExpiredToken    = errors.New("token is expired")
	ErrAuthInvalidEmail    = errors.New("invalid email")
	ErrAuthInvalidPassword = errors.New("invalid password")
	ErrAuthCreateToken     = errors.New("failed to create token")
	ErrAuthHashPassword    = errors.New("failed to hash password")
)
