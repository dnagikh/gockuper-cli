package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var defaults = map[string]interface{}{
	"DB_TYPE":            "postgres",
	"STORAGE_TYPE":       "dropbox",
	"LOG_NAME":           "./gockuper.log",
	"LOG_MAX_SIZE":       50, // MB
	"LOG_MAX_BACKUPS":    3,
	"LOG_MAX_AGE":        90, // Days
	"LOG_COMPRESSION":    true,
	"LOG_LEVEL":          "info",
	"LOG_TARGET":         "stdout",
	"STORAGE_FILE_PATH":  "/",
	"COMPRESS":           "none",
	"MAX_BACKUPS":        10,
	"DROPBOX_TOKEN_FILE": "./",
}

func LoadConfig() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(exePath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	for key, value := range defaults {
		if !viper.IsSet(key) {
			viper.Set(key, value)
		}
	}

	return nil
}
