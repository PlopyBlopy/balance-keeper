package testdata

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Драйвер для БД
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Источник "файлы"
)

func InitMigrate(folderPath string, connString string) error {
	rawPath := filepath.ToSlash(folderPath)
	ulrPath := fmt.Sprintf("file://%s", rawPath)

	m, err := migrate.New(ulrPath, connString)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	if dbErr != nil {
		return dbErr
	}

	return nil
}
