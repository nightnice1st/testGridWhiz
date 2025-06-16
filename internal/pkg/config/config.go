package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI        string
	JWTSecret         string
	JWTExpiry         time.Duration
	AuthServicePort   string
	UserServicePort   string
	RateLimitAttempts int
	RateLimitWindow   time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	jwtExpiry, _ := time.ParseDuration(os.Getenv("JWT_EXPIRY"))
	rateLimitWindow, _ := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))

	return &Config{
		MongoDBURI:        os.Getenv("MONGODB_URI"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiry:         jwtExpiry,
		AuthServicePort:   os.Getenv("AUTH_SERVICE_PORT"),
		UserServicePort:   os.Getenv("USER_SERVICE_PORT"),
		RateLimitAttempts: 5,
		RateLimitWindow:   rateLimitWindow,
	}
}
