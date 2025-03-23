package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dnagikh/gockuper-cli/internal/logger"
	"github.com/spf13/viper"
)

var tokenFileDefaultPath = "./token.json"

type DropboxToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type DropboxTokenProvider struct {
	clientID     string
	clientSecret string
	tokenPath    string
	current      DropboxToken
	lock         sync.RWMutex
}

func NewDropboxTokenProvider(clientID, clientSecret string) (*DropboxTokenProvider, error) {
	tokenPath := viper.GetString("DROPBOX_TOKEN_FILE")
	if tokenPath == "" {
		tokenPath = tokenFileDefaultPath
	}

	token, err := loadTokenFromFile(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load dropbox token from file: %w", err)
	}

	return &DropboxTokenProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenPath:    tokenPath,
		current:      *token,
	}, nil
}

func (p *DropboxTokenProvider) AccessToken() string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.current.AccessToken
}

func (p *DropboxTokenProvider) StartBackgroundRefresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			if err := p.refreshTokenIfNeeded(); err != nil {
				logger.Log.Error("failed to refresh token via Dropbox API", slog.String("error", err.Error()))
			}

			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (p *DropboxTokenProvider) refreshTokenIfNeeded() error {
	p.lock.RLock()
	expiresAt := p.current.ExpiresAt
	p.lock.RUnlock()

	expiresSoon := time.Until(expiresAt) < (15 * time.Minute)

	if !expiresSoon {
		return nil
	}

	return p.refreshToken()
}

func (p *DropboxTokenProvider) refreshToken() error {
	p.lock.Lock()
	defer p.lock.Unlock()

	form := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s",
		p.current.RefreshToken, p.clientID, p.clientSecret)

	req, err := http.NewRequest("POST", "https://api.dropboxapi.com/oauth2/token", io.NopCloser(strings.NewReader(form)))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("invalid token response: %d %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	if result.RefreshToken != "" {
		p.current.RefreshToken = result.RefreshToken
	}
	p.current.AccessToken = result.AccessToken
	p.current.ExpiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)

	if err := saveTokenToFile(p.tokenPath, &p.current); err != nil {
		return fmt.Errorf("failed to save token to file: %w", err)
	}

	return nil
}

func loadTokenFromFile(path string) (*DropboxToken, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}
	defer f.Close()

	t := DropboxToken{}
	if err := json.NewDecoder(f).Decode(&t); err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}
	return &t, nil
}

func saveTokenToFile(path string, token *DropboxToken) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
