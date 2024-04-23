package initializers

import (
	"fmt"
	// "log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env variables (log)")
		// log.Fatal("Error loading .env variables (fatal)")
	}
}
