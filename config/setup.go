package config

import (
	"os"
	"strconv"
)

type Config struct {
	EmailFrom      string
	EmailHost      string
	EmailPort      int
	EmailUsername  string
	EmailPassword  string
	DisplayName    string
	Hostname       string
	StaticHostname string
	RedisURI       string
	DBPath         string
	SecretKey      string
}

var AppConfig Config

func Init() {
	AppConfig.SecretKey = getEnvOrPanic("SECRET_KEY")

	AppConfig.EmailFrom = getEnvOrPanic("EMAIL_FROM")
	AppConfig.EmailHost = getEnvOrPanic("EMAIL_HOST")

	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		panic("Invalid EMAIL_PORT")
	}
	AppConfig.EmailPort = port
	AppConfig.EmailUsername = getEnvOrPanic("EMAIL_USERNAME")
	AppConfig.EmailPassword = getEnvOrPanic("EMAIL_PASSWORD")

	AppConfig.DisplayName = getEnvOrPanic("DISPLAY_NAME")
	AppConfig.Hostname = getEnvOrPanic("HOSTNAME")

	AppConfig.StaticHostname = os.Getenv("STATIC_HOSTNAME")

	AppConfig.RedisURI = getEnvOrPanic("REDIS_URI")
	AppConfig.DBPath = getEnvOrPanic("DB_PATH")
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Missing environment variable: " + key)
	}
	return value
}
