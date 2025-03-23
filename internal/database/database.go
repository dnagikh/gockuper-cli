package database

import (
	"fmt"
	"io"

	"github.com/spf13/viper"
)

type Database interface {
	Dump() (io.Reader, error)
	Version() (string, error)
}

func NewDatabase() (Database, error) {
	switch viper.GetString("DB_TYPE") {
	case "postgres":
		return NewPostgresDB(), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", viper.GetString("DB_TYPE"))
	}
}
