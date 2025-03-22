package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

func InitLogger() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %v", err)
	}

	var level slog.Level
	switch viper.GetString("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	default:
		level = slog.LevelError
	}

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}
	var handler *slog.JSONHandler

	switch viper.GetString("LOG_TARGET") {
	case "file":
		logFilePath := filepath.Join(filepath.Dir(exePath), viper.GetString("LOG_NAME"))

		logFile := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    viper.GetInt("LOG_MAX_SIZE"),
			MaxBackups: viper.GetInt("LOG_MAX_BACKUPS"),
			MaxAge:     viper.GetInt("LOG_MAX_AGE"),
			Compress:   viper.GetBool("LOG_COMPRESSION"),
		}

		handler = slog.NewJSONHandler(logFile, handlerOpts)
	case "stdout", "":
		handler = slog.NewJSONHandler(os.Stdout, handlerOpts)
	default:
		return fmt.Errorf("unsupported LOG_TARGET: %s", viper.GetString("LOG_TARGET"))
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)

	Log.Info("Logger initialized",
		"target", viper.GetString("LOG_TARGET"),
		"level", level.String(),
	)

	return nil
}
