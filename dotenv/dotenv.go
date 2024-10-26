package dotenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Redis struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Env struct {
	Redis Redis
}

func LoadEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUser := os.Getenv("REDIS_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	return &Env{
		Redis: Redis{
			Host:     redisHost,
			Port:     redisPort,
			User:     redisUser,
			Password: redisPassword,
		},
	}
}
