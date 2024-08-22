package cache

import (
	"context"

	"coinex-api/v0/cookies"
	"coinex-api/v0/db"
	"coinex-api/v0/pkg/views"
)

var Repository views.Cache

func Init(storeType, databaseURL string, ctx context.Context) views.Cache {
	var cache views.Cache
	switch storeType {
	case "DB":
		db.Repository.Init(ctx, databaseURL)
		cache.Store = &db.Repository
	default:
		cache.Store = &cookies.Repository
	}
	return cache
}

func Close(cache *views.Cache) {
	defer cache.Store.Close()
}
