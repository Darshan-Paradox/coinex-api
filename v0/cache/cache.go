package cache

import (
	"context"
	"os"

	"coinex-api/v0/cookies"
	"coinex-api/v0/db"
	"coinex-api/v0/pkg/views"
)

var Repository Cache

type ICache interface {
	GetAllCoin() (views.CoinList, error)
	GetCoin(coinCode string) (views.Coin, error)
	GetPrice(coinCode string) (views.PriceResponse, error)
	GetCoinInCurrency(coinCode, currency string) (views.Coin, error)
	GetPriceIn(coinCode, currency string) (views.Price, error)

	SetCoin(coin views.Coin) error
	SetPrice(price views.Price) error

	Close()
}

// Factory Generator
type Cache struct {
	Store ICache
}

func (cache *Cache) Init(ctx context.Context) {
	switch cacheEnv := os.Getenv("CACHE"); cacheEnv {
	case "DB":
		db.Repository.Init(ctx, os.Getenv("DATABASE_URL"))
		cache.Store = &db.Repository
	default:
		cache.Store = &cookies.Repository
	}
}
