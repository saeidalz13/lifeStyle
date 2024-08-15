package handlers_test

import (
	"log"
	"os"
	"testing"

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

	ah := handlers.NewAuthHandler(tm)
	authHandler = ah

	os.Exit(m.Run())
}
