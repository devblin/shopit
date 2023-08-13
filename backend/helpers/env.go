package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const PRODUCTION string = "production"
const DEV string = "dev"

func init() {
	if os.Getenv("ENV") != PRODUCTION {
		var err = godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Print("LOADED ENVIRONMENT VARIABLES...")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
