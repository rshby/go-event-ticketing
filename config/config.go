package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// LoadConfig loads config from env file
func LoadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}
}
