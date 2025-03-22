package cmd

import (
	"github.com/dnagikh/gockuper-cli/internal/backup"
	"github.com/dnagikh/gockuper-cli/internal/logger"
	"github.com/spf13/cobra"
	"log/slog"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create database backup and upload it to the cloud",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Log.Info("starting backup...")

		err := backup.CreateBackup()
		if err != nil {
			logger.Log.Error("failed while creating backup", slog.String("desc", err.Error()))
			return
		}

		logger.Log.Info("backup successfully created")
	},
}
