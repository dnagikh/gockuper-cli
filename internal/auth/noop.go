package auth

import "context"

type NoopRefresher struct{}

func (NoopRefresher) StartBackgroundRefresh(_ context.Context) {
}
