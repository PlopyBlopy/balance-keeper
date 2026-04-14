package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PlopyBlopy/balance-keeper-service/internal/shared"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	IsProd bool
	// HTTP
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
	// DB
	DBConnString string `env:"APPDB_DBConnString"`
	MinConns     int32  `env:"APPDB_MINCONNS"`
	MaxConns     int32  `env:"APPDB_MAXCONNS"`
}

// Для определения окружения пытается получить значение из флага, если empty то из environments, если empty то IsProd = false
// Имя файла .env должно соответствовать переданному окружению, example: .env.dev
func NewConfig() (*Config, error) {
	c := &Config{}

	envFlag := flag.String("env", "dev", "Environment: dev|prod|development|production")

	flag.Parse()

	if *envFlag != "" {
		if *envFlag == "prod" || *envFlag == "production" {
			c.IsProd = true
		}
	} else {
		curenv := os.Getenv("env")
		if curenv != "" {
			if curenv == "prod" || curenv == "production" {
				c.IsProd = true
			}
		}
	}

	root, err := shared.FindProjectRoot()
	if err != nil {
		// Если корень не был найден, не считается ошибкой так как приложение может быть бинарным файлом
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return nil, err
	} else {
		envFile := filepath.Join(root, fmt.Sprintf(".env.%s", *envFlag))

		err := godotenv.Load(envFile)
		// Если используется иная загрузка env например из docker, а не файл в корне
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	err = env.Parse(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
