package testdata

import (
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type TestConfig struct {
	// HTTP
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
	// DB - Testcontainer
	Image        string `env:"DB_IMAGE"`
	Database     string `env:"DB_DATABASE"`
	Username     string `env:"DB_USERNAME"`
	Password     string `env:"DB_PASSWORD"`
	MinConns     int32  `env:"DB_MINCONNS"`
	MaxConns     int32  `env:"DB_MAXCONNS"`
	ShapshotName string `env:"SNAPSHOTNAME"`
}

// load .env from cur folder
func NewTestConfig() (*TestConfig, error) {
	c := &TestConfig{}

	curdir := GetCurDirPath()

	err := godotenv.Load(filepath.Join(curdir, ".env"))
	if err != nil {
		return nil, err
	}

	err = env.Parse(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
