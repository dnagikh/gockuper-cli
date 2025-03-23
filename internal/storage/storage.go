package storage

import (
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/dnagikh/gockuper-cli/internal/auth"
	"github.com/spf13/viper"
)

type Storage interface {
	Upload(file io.Reader, fileName string) error
	ListFiles(folder string) ([]StoredFile, error)
	Delete(filename string) error
}

func NewStorage() (Storage, error) {
	switch viper.GetString("STORAGE_TYPE") {
	case "dropbox":
		clientId := viper.GetString("DROPBOX_CLIENT_ID")
		clientSecret := viper.GetString("DROPBOX_CLIENT_SECRET")
		provider, err := auth.NewDropboxTokenProvider(clientId, clientSecret)
		if err != nil {
			return nil, fmt.Errorf("could not create Dropbox token provider: %w", err)
		}

		return NewDropbox(provider), nil
	case "file":
		return NewFileStorage(), nil
	default:
		return nil, fmt.Errorf("unsupported storage type %s", viper.GetString("STORAGE_TYPE"))
	}
}

type ByteCounter struct {
	Total int64
}

func (bc *ByteCounter) Write(p []byte) (int, error) {
	n := len(p)
	bc.Total += int64(n)
	return n, nil
}

type StoredFile struct {
	Name      string
	Timestamp time.Time
}

var nameDateRegexp = regexp.MustCompile(`^dump_\d+\.\d+_(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})`)

func parseTimeFromFilename(name string) (time.Time, error) {
	matches := nameDateRegexp.FindStringSubmatch(name)
	if len(matches) < 2 {
		return time.Time{}, fmt.Errorf("filename does not match expected pattern: %s", name)
	}
	return time.Parse("2006-01-02_15-04-05", matches[1])
}
