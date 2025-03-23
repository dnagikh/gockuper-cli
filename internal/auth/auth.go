package auth

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

type TokenRefresher interface {
	StartBackgroundRefresh(ctx context.Context)
}

func NewTokenRefresher() (TokenRefresher, error) {
	var refresher TokenRefresher
	storageType := viper.GetString("STORAGE_TYPE")

	switch storageType {
	case "dropbox":
		provider, err := NewDropboxTokenProvider(viper.GetString("DROPBOX_CLIENT_ID"), viper.GetString("DROPBOX_CLIENT_SECRET"))
		if err != nil {
			return nil, fmt.Errorf("failed while creating dropbox token provider: %w", err)
		}
		refresher = provider
	case "file":
		refresher = NoopRefresher{}
	default:
		return nil, fmt.Errorf("unsupported storage type")
	}

	return refresher, nil
}
