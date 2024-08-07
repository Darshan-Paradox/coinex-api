package main

import (
	"context"
	"os"

	"coinex-api/v0/cache"
	"coinex-api/v0/logger"
	"coinex-api/v0/pkg/routes"
)

var repo = &cache.Repository
var router = &routes.RouterInst
var log = &logger.LoggerInst

func main() {

	log.Init()
	defer log.Close()

	repo = cache.Init(os.Getenv("CACHE_IMPLEMENTATION"), os.Getenv("DATABASE_URL"), context.Background())
	defer cache.Close(repo)

	router.Init()
	router.Run(os.Getenv("PORT"))
}
