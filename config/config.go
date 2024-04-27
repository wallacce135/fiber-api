package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading.env file")
		os.Exit(1)
	}
	return os.Getenv(key)
}
