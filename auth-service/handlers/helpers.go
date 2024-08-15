package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/saeidalz13/lifestyle/auth-service/models"
)

func writeApiErrResp(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	errResp := models.NewAuthApiRespWithErr[models.NoPayload](err)
	json.NewEncoder(w).Encode(errResp)
}
