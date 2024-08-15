package main

import (
	"auth-service/internal/apiutils"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("STAGE") != apiutils.ApiStageProd {
		if err := godotenv.Load("./config/.env.local"); err != nil {
			log.Fatalln(err)
		}
	}
	port := os.Getenv("PORT")

	server := &http.Server{
		Addr: "127.0.0.1" + ":" + port,
	}

	log.Printf("Listening to %s...\n", server.Addr)
	log.Fatalln(server.ListenAndServe())
}
