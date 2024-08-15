package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/saeidalz13/lifestyle/auth-service/internal/autherr"
	"github.com/saeidalz13/lifestyle/auth-service/models"
	"github.com/saeidalz13/lifestyle/auth-service/routes"
)

type AuthTestCase[T any] struct {
	name               string
	reqBody            models.ReqAuth
	expectedStatusCode int
	expectedResp       models.AuthApiResp[T]
}

func TestSignup(t *testing.T) {
	testCases := []AuthTestCase[any]{
		{
			name:               "should signup valid email and password",
			reqBody:            models.ReqAuth{Email: validEmail, Password: []byte(validPassword)},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "fail signup invalid email",
			reqBody:            models.ReqAuth{Email: invalidEmail, Password: []byte(validPassword)},
			expectedStatusCode: http.StatusBadRequest,
			expectedResp:       models.NewAuthApiRespWithErr[any](autherr.ErrAuthInvalidEmail),
		},
		{
			name:               "fail signup invalid password",
			reqBody:            models.ReqAuth{Email: validEmail, Password: []byte(invalidPassword)},
			expectedStatusCode: http.StatusBadRequest,
			expectedResp:       models.NewAuthApiRespWithErr[any](autherr.ErrAuthShortPassword),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			r := httptest.NewRequest(http.MethodPost, routes.Signup, bytes.NewBuffer(body))
			r.Header.Set("Content-Type", validContentType)
			w := httptest.NewRecorder()

			authHandler.ServeHTTP(w, r)

			respStatusCode := w.Result().StatusCode

			if respStatusCode != test.expectedStatusCode {
				t.Fatalf("expected status: %d\tgot: %d", test.expectedStatusCode, respStatusCode)
			}

			if respStatusCode != http.StatusCreated {
				respBytes, err := io.ReadAll(w.Result().Body)
				if err != nil {
					t.Fatal(err)
				}

				var resp models.AuthApiResp[any]
				if err := json.Unmarshal(respBytes, &resp); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(test.expectedResp, resp) {
					t.Fatalf("expected resp: %+v\ngot: %+v\n", test.expectedResp, resp)
				}
			}
		})
	}
}
