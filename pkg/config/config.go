// Package config manages the configuration
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Addr        string
	DatabaseURI string
	Secret      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	configuration := Config{
		Port:        check("SERVER_PORT", "8888"),
		Addr:        check("SERVER_ADDRESS", "0.0.0.0"),
		DatabaseURI: check("DATABASE_URI", ""),
		Secret:      check("SECRET", "Secret"),
	}
	return &configuration
}

func check(parameter, defaultValue string) string {
	if os.Getenv(parameter) == "" {
		return defaultValue
	}
	return os.Getenv(parameter)
}
