package backup

import (
	"fmt"
	"github.com/dnagikh/gockuper-cli/internal/compress"
	"github.com/dnagikh/gockuper-cli/internal/database"
	"github.com/dnagikh/gockuper-cli/internal/logger"
	"github.com/dnagikh/gockuper-cli/internal/storage"
	"github.com/spf13/viper"
	"log/slog"
	"sort"
	"time"
)

func CreateBackup() error {
	db, err := database.NewDatabase()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	dbVersion, err := db.Version()
	if err != nil {
		return fmt.Errorf("could not get db version: %w", err)
	}

	file, err := db.Dump()
	if err != nil {
		return fmt.Errorf("could not dump database: %w", err)
	}

	compressor, err := compress.FromString(viper.GetString("COMPRESS"))
	if err != nil {
		return fmt.Errorf("could not compress database: %w", err)
	}

	compressedFile, err := compress.Compress(file, compressor)
	if err != nil {
		return fmt.Errorf("could not compress database: %w", err)
	}

	strg, err := storage.NewStorage()
	if err != nil {
		return fmt.Errorf("could not connect to storage: %w", err)
	}

	filename := fmt.Sprintf("dump_%s_%s.%s", dbVersion, time.Now().Format("2006-01-02_15-04-05"), compressor.Extension())
	if err := strg.Upload(compressedFile, filename); err != nil {
		return fmt.Errorf("could not upload database: %w", err)
	}

	maxBackups := viper.GetInt("MAX_BACKUPS")
	folder := viper.GetString("STORAGE_FILE_PATH")

	if err := CleanupOldBackups(strg, folder, maxBackups); err != nil {
		return fmt.Errorf("could not cleanup old backups: %w", err)
	}

	return nil
}

func CleanupOldBackups(strg storage.Storage, folder string, max int) error {
	files, err := strg.ListFiles(folder)
	if err != nil {
		return fmt.Errorf("could not get list files in dir %s: %w", folder, err)
	}

	logger.Log.Info("found files", slog.Int("count", len(files)))

	if len(files) < max {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp.After(files[j].Timestamp)
	})

	toDelete := files[max:]

	for _, file := range toDelete {
		err := strg.Delete(file.Name)
		if err != nil {
			logger.Log.Error("could not delete file", slog.String("filename", file.Name), slog.String("error", err.Error()))
			continue
		}
	}

	return nil
}
