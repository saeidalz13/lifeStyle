package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/saeidalz13/lifestyle/auth-service/handlers"
	"github.com/saeidalz13/lifestyle/auth-service/internal/apiutils"
	"github.com/saeidalz13/lifestyle/auth-service/token"
)

func main() {
	log.SetFlags(log.Lshortfile)
	if os.Getenv("STAGE") != apiutils.ApiStageProd {
		if err := godotenv.Load("./config/.env.local"); err != nil {
			log.Fatalln(err)
		}
	}
	port := os.Getenv("PORT")

	tokenManager, err := token.BuildPasetoTokenManager()
	if err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{
		Addr:    "127.0.0.1" + ":" + port,
		Handler: handlers.NewAuthHandler(tokenManager),
	}

	log.Printf("Listening to %s...\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
}
