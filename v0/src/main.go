package main

import (
	"os"
	"context"

	"coinex-api/v0/logger"
	"coinex-api/v0/pkg/routes"
	"coinex-api/v0/internal/cache"
)

var repo = &cache.Repository
var router = &routes.RouterInst
var log = &logger.LoggerInst

func main() {

	log.Init()
	defer log.Close()

	repo.Init(context.Background())
	defer repo.Store.Close()

	router.Init()
	router.Run(os.Getenv("PORT"))
}
