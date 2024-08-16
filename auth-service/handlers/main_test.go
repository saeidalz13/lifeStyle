package handlers_test

import (
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/saeidalz13/lifestyle/auth-service/handlers"
	"github.com/saeidalz13/lifestyle/auth-service/token"
)

const (
	validEmail       = "saeid23@hotmail.com"
	validPassword    = "validPass@3525"
	validContentType = "application/json"

	invalidEmail    = "efgag"
	invalidPassword = "egag"
)

var (
	validToken  string
	authHandler *handlers.AuthHandler

	testMock sqlmock.Sqlmock
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func TestMain(m *testing.M) {
	tm, err := token.BuildPasetoTokenManager()
	checkErr(err)

	vt, err := tm.CreateToken(validEmail, token.PasetoCookieDuration)
	checkErr(err)
	validToken = vt

	db, mock, err := sqlmock.New()
	checkErr(err)
	testMock = mock

	ah := handlers.NewAuthHandler(tm, db)
	authHandler = ah

	os.Exit(m.Run())
}
