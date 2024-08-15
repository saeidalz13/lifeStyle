package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/saeidalz13/lifestyle/auth-service/internal/apierr"
	"github.com/saeidalz13/lifestyle/auth-service/internal/autherr"
	"github.com/saeidalz13/lifestyle/auth-service/models"
	"github.com/saeidalz13/lifestyle/auth-service/token"
	"golang.org/x/crypto/bcrypt"
)

const (
	ActionLogin  string = "login"
	ActionSignup string = "signup"
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
	if err := json.Unmarshal(bodyBytes, &reqAuth); err != nil {
		writeApiErrResp(w, http.StatusBadRequest, apierr.ErrApiUnmarshalBody)
		return
	}

	var statusCode int
	switch r.URL.Query().Get("action") {
	case ActionLogin:
		if !a.processLogin(w, &reqAuth) {
			return
		}
		statusCode = http.StatusOK

	case ActionSignup:
		if !a.processSignUp(w, &reqAuth) {
			return
		}
		statusCode = http.StatusCreated

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

func isEmailValid(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return autherr.ErrAuthShortPassword
	}

	// The regex pattern `[!@#$%^&*(),.?":{}|<>]` matches common special characters
	regex := `[!@#$%^&*(),.?":{}|<>]`
	specialCharRe := regexp.MustCompile(regex)

	// Patter for digit existence
	digitRegex := `[0-9]`
	digitRe := regexp.MustCompile(digitRegex)

	// Check if the password matches the regex pattern
	if specialCharRe.MatchString(password) && digitRe.MatchString(password) {
        return autherr.ErrAuthInvalidPassword
    }

    return nil
}
