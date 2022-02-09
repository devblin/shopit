package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var PRODUCTION string = "production"

func init() {
	if os.Getenv("ENV") != PRODUCTION {
		var err = godotenv.Load("./env/.env")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("LOADED ENVIRONMENT VARIABLES")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
