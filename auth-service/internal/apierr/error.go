package apierr

import "errors"

var (
	ErrApiReadBody            = errors.New("failed to read body bytes")
	ErrApiUnmarshalBody       = errors.New("failed to unmarshal body")

)
