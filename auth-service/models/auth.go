package models

type ReqAuth struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
