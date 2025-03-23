package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type PostgresDB struct{}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

func (*PostgresDB) Dump() (io.Reader, error) {
	pr, pw := io.Pipe()

	cmd := exec.Command(
		"pg_dump",
		"-F", "c",
		"-h", viper.GetString("DB_HOST"),
		"-p", viper.GetString("DB_PORT"),
		"-U", viper.GetString("DB_USER"),
		"-d", viper.GetString("DB_NAME"),
	)

	cmd.Env = append(os.Environ(), "PGPASSWORD="+viper.GetString("DB_PASSWORD"))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get stdout pipe: %v", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("pg_dump start failed: %w", err)
	}

	go func() {
		defer pw.Close()

		if _, err := io.Copy(pw, stdout); err != nil {
			pw.CloseWithError(fmt.Errorf("pg_dump stdout failed: %w", err))
			return
		}
	}()

	return pr, nil
}

func (*PostgresDB) Version() (string, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "", fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SHOW server_version").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to get server version: %w", err)
	}

	parts := strings.Fields(version)
	if len(parts) > 0 {
		version = parts[0]
	}

	return version, nil
}
