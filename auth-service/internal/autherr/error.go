package autherr

import "errors"

var (
    ErrAuthExpiredToken = errors.New("token is expired")
)