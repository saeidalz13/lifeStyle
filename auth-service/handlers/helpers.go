package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/saeidalz13/lifestyle/auth-service/internal/autherr"
	"github.com/saeidalz13/lifestyle/auth-service/models"
)

func writeApiErrResp(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	errResp := models.NewAuthApiRespWithErr[models.NoPayload](err)
	json.NewEncoder(w).Encode(errResp)
}

func extractAuthBody(bodyBytes []byte) (models.ReqAuth, error) {
	var reqAuth models.ReqAuth
	if err := json.Unmarshal(bodyBytes, &reqAuth); err != nil {
		return reqAuth, err
	}

	return reqAuth, nil
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
