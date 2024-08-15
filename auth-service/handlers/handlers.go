package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/saeidalz13/lifestyle/auth-service/models"
	"github.com/saeidalz13/lifestyle/auth-service/token"
)

type AuthHandler struct {
	tm token.TokenManger
}

func NewAuthHandler(tm token.TokenManger) AuthHandler {
	return AuthHandler{tm: tm}
}

func (a AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pasetoCookie, err := r.Cookie(token.PasetoCookieName)
	if err != nil {
        writeApiErrResp(w, http.StatusUnauthorized, err)
	}

    a.tm.VerifyToken()
}

func writeApiErrResp(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
    errResp := models.NewAuthApiRespWithErr[models.NoPayload](err)
	json.NewEncoder(w).Encode(errResp)
}
