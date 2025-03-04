package config

import (
	"embed"
	"os"

	"github.com/joho/godotenv"
)

//go:embed ".env"
var envFile embed.FS

func LoadEnv() error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return err
		}
		return nil
	}

	data, err := envFile.ReadFile(".env")
	if err != nil {
		return err
	}

	envMap, err := godotenv.Unmarshal(string(data))
	if err != nil {
		return err
	}

	for key, value := range envMap {
		os.Setenv(key, value)
	}

	return nil
}
