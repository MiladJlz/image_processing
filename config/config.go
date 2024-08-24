package config

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	mw "image_processing/api/middleware"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host             string
	Port             string
	PostgresPassword string
	PostgresUser     string
	PostgresDBName   string
	PostgresPort     string
	SSLMode          string
	RedisDB          string
	RedisPass        string
	RedisPort        string
	JwtSecret        string
}

func GetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	return &Config{
		Host:             os.Getenv("DB_HOST"),
		Port:             os.Getenv("PORT"),
		PostgresPassword: os.Getenv("DB_PASS"),
		PostgresUser:     os.Getenv("DB_USER"),
		PostgresDBName:   os.Getenv("DB_NAME"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		SSLMode:          os.Getenv("SSLMODE"),
		RedisDB:          os.Getenv("REDIS_DB"),
		RedisPass:        os.Getenv("REDIS_PASS"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		JwtSecret:        os.Getenv("SECRET"),
	}

}
func (c *Config) GetPostgresConn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDBName, c.SSLMode,
	)
}

func (c *Config) GetRedisOptions() *redis.Options {
	db, _ := strconv.Atoi(c.RedisDB)
	return &redis.Options{
		Addr:     fmt.Sprint(c.Host + ":" + c.RedisPort),
		Password: c.RedisPass,
		DB:       db,
	}
}

func (c *Config) GetJwtConfig() echojwt.Config {
	return echojwt.Config{ContinueOnIgnoredError: true,
		ErrorHandler: mw.JWTAuthentication,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(mw.JwtCustomClaims)
		},
		SigningKey: []byte(c.JwtSecret),
	}
}
