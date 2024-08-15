package autherr

import "errors"

var (
	ErrAuthExpiredToken    = errors.New("token is expired")
	ErrAuthInvalidEmail    = errors.New("invalid email")
	ErrAuthShortPassword   = errors.New("password must be a minimum of 8 characters")
	ErrAuthInvalidPassword = errors.New("password must contain at least one uppercase letter and one digit")
	ErrAuthCreateToken     = errors.New("failed to create token")
	ErrAuthHashPassword    = errors.New("failed to hash password")
)
