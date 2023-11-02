package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitEnvironment() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func GetValue(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		panic("Environment variable " + key + " does not exist")
	}

	return value
}
