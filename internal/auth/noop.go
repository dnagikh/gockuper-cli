package auth

import "context"

type NoopRefresher struct{}

func (NoopRefresher) StartBackgroundRefresh(ctx context.Context) {
}
