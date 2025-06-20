package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr   string `mapstructure:"SERVER_ADDR"`
	DBUrl        string `mapstructure:"DATABASE_URL"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadConfig() (config Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Println("could not load .env file")
	}

	config.ServerAddr = os.Getenv("SERVER_ADDR")
	config.DBUrl = os.Getenv("DATABASE_URL")
	config.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	return
}
