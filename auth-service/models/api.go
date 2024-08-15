package models

type NoPayload bool

type AuthApiResp[T any] struct {
	Payload T      `json:"payload"`
	Err     string `json:"paylod"`
}

func NewAuthApiRespWithPayload[T any](p T) AuthApiResp[T] {
	return AuthApiResp[T]{Payload: p}
}

func NewAuthApiRespWithErr[T any](err error) AuthApiResp[T] {
	return AuthApiResp[T]{Err: err.Error()}
}
