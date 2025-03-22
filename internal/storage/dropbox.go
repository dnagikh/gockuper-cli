package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dnagikh/gockuper-cli/internal/auth"
	"github.com/dnagikh/gockuper-cli/internal/logger"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	listFilesUrl = "https://api.dropboxapi.com/2/files/list_folder"
	uploadUrl    = "https://content.dropboxapi.com/2/files/upload"
	deleteUrl    = "https://api.dropboxapi.com/2/files/delete_v2"
)

type Dropbox struct {
	tokenProvider *auth.DropboxTokenProvider
}

func NewDropbox(provider *auth.DropboxTokenProvider) *Dropbox {
	return &Dropbox{
		tokenProvider: provider,
	}
}

func (d *Dropbox) Upload(file io.Reader, filename string) error {
	counter := &ByteCounter{}
	tee := io.TeeReader(file, counter)

	start := time.Now()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(tee)
	if err != nil {
		return fmt.Errorf("could not read from file: %w", err)
	}

	req, err := http.NewRequest("POST", uploadUrl, buf)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	token := d.tokenProvider.AccessToken()
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Dropbox-API-Arg", fmt.Sprintf(`{"path": "%s%s"}`, viper.GetString("STORAGE_FILE_PATH"), filename))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody := new(bytes.Buffer)
		_, _ = errBody.ReadFrom(resp.Body)

		return fmt.Errorf(
			"dropbox upload failed: status=%d, error=%s",
			resp.StatusCode,
			errBody.String(),
		)
	}

	duration := time.Since(start)
	sizeMB := float64(counter.Total) / 1024.0 / 1024.0

	logger.Log.Info("dropbox upload success", slog.String("size", fmt.Sprintf("%.2f MB", sizeMB)), slog.String("duration", fmt.Sprintf("%.2f s", duration.Seconds())))

	return nil
}

func (d *Dropbox) ListFiles(folder string) ([]StoredFile, error) {
	if folder == "/" {
		folder = ""
	}
	reqBody := fmt.Sprintf(`{"path": "%s"}`, folder)
	resp, err := doDropboxRequest(listFilesUrl, reqBody, d.tokenProvider.AccessToken())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed struct {
		Entries []struct {
			Name string `json:"name"`
		} `json:"entries"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	files := make([]StoredFile, 0, len(parsed.Entries))
	for _, entry := range parsed.Entries {
		t, err := parseTimeFromFilename(entry.Name)
		if err != nil {
			continue
		}

		files = append(files, StoredFile{
			Name:      entry.Name,
			Timestamp: t,
		})
	}

	return files, nil
}

func (d *Dropbox) Delete(filename string) error {
	fullPath := fmt.Sprintf(`{"path": "%s%s"}`, viper.GetString("STORAGE_FILE_PATH"), filename)
	resp, err := doDropboxRequest(deleteUrl, fullPath, d.tokenProvider.AccessToken())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func doDropboxRequest(url, body, token string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dropbox request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("dropbox API error (%d): %s", resp.StatusCode, string(data))
	}

	return resp, nil
}
