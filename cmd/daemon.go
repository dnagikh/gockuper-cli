package cmd

import (
	"context"
	"github.com/dnagikh/gockuper-cli/internal/auth"
	"github.com/dnagikh/gockuper-cli/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Cloud service authentication / refreshing tokens",
	RunE: func(cmd *cobra.Command, args []string) error {
		refresher, err := auth.NewTokenRefresher()
		if err != nil {
			logger.Log.Error("failed to initialize token refresher with error", slog.String("error", err.Error()))
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		refresher.StartBackgroundRefresh(ctx)
		logger.Log.Info("Token refresher started for storage", slog.String("type", viper.GetString("STORAGE_TYPE")))

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Log.Info("Daemon stopped")

		return nil
	},
}
