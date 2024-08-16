package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/saeidalz13/lifestyle/auth-service/internal/apierr"
	"github.com/saeidalz13/lifestyle/auth-service/internal/autherr"
	"github.com/saeidalz13/lifestyle/auth-service/models"
	"github.com/saeidalz13/lifestyle/auth-service/routes"
	"github.com/saeidalz13/lifestyle/auth-service/token"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	tokenManager token.TokenManger
}

func NewAuthHandler(ptm token.TokenManger) *AuthHandler {
	return &AuthHandler{tokenManager: ptm}
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		writeApiErrResp(w, http.StatusInternalServerError, apierr.ErrApiReadBody)
		return
	}

	var reqAuth models.ReqAuth
	var statusCode int
	switch r.URL.Path {
	case routes.Login:
		reqAuth, err = extractAuthBody(bodyBytes)
		if err != nil {
			writeApiErrResp(w, http.StatusBadRequest, apierr.ErrApiUnmarshalBody)
		}
		if !a.processLogin(w, &reqAuth) {
			return
		}
		statusCode = http.StatusOK

	case routes.Signup:
		reqAuth, err = extractAuthBody(bodyBytes)
		if err != nil {
			writeApiErrResp(w, http.StatusBadRequest, apierr.ErrApiUnmarshalBody)
		}
		if !a.processSignUp(w, &reqAuth) {
			return
		}
		statusCode = http.StatusCreated

	case routes.TokenAuth:
		a.authenticateReqByCookie(w, r)
		return

	default:
		http.Error(w, "invalid action of auth", http.StatusBadRequest)
		return
	}

	pasetoToken, err := a.tokenManager.CreateToken(reqAuth.Email, token.PasetoCookieDuration)
	if err != nil {
		log.Println(err)
		writeApiErrResp(w, http.StatusInternalServerError, autherr.ErrAuthCreateToken)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     token.PasetoCookieName,
		Value:    pasetoToken,
		Expires:  time.Now().Add(token.PasetoCookieDuration),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		// Secure:   true,
	})

	w.WriteHeader(statusCode)
}

func (a *AuthHandler) processSignUp(w http.ResponseWriter, reqAuth *models.ReqAuth) bool {
	if !isEmailValid(reqAuth.Email) {
		writeApiErrResp(w, http.StatusBadRequest, autherr.ErrAuthInvalidEmail)
		return false
	}

	if err := validatePassword(string(reqAuth.Password)); err != nil {
		writeApiErrResp(w, http.StatusBadRequest, err)
		return false
	}

	_, err := bcrypt.GenerateFromPassword(reqAuth.Password, bcrypt.DefaultCost)
	if err != nil {
		writeApiErrResp(w, http.StatusInternalServerError, autherr.ErrAuthHashPassword)
		return false
	}

	// TODO: Add to db

	return true
}

func (a *AuthHandler) processLogin(w http.ResponseWriter, reqAuth *models.ReqAuth) bool {
	// TODO: Some db to extract the password. for now random

	randomHashed := []byte("somepass")
	if err := bcrypt.CompareHashAndPassword(randomHashed, reqAuth.Password); err != nil {
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

func (a *AuthHandler) authenticateReqByCookie(w http.ResponseWriter, r *http.Request) {
	pasetoCookie, err := r.Cookie(token.PasetoCookieName)
	if err != nil {
		writeApiErrResp(w, http.StatusUnauthorized, err)
		return
	}

	pasetoPayload, err := a.tokenManager.VerifyToken(pasetoCookie.Value)
	if err != nil {
		writeApiErrResp(w, http.StatusUnauthorized, err)
		return
	}

	json.NewEncoder(w).Encode(models.RespTokenAuth{Email: pasetoPayload.Email})
}
