package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	GRPCServerPort  string
	LoggerLevel     string
	DebugServerPort string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		GRPCServerPort:  os.Getenv("GRPC_SERVER_PORT"),
		LoggerLevel:     os.Getenv("LOGGER_LEVEL"),
		DebugServerPort: os.Getenv("DEBUG_SERVER_PORT"),
	}
}
