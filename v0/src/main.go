package main

import (
	"os"
	"context"

	"coinex-api/v0/db"
	"coinex-api/v0/logger"
	"coinex-api/v0/pkg/routes"
)

var repo = &db.Repository
var router = &routes.RouterInst
var log = &logger.LoggerInst

func main() {

	log.Init()
	defer log.Close()

	repo.InitDB(context.Background(), os.Getenv("DATABASE_URL"))
	defer repo.Close()

	router.Init()
	router.Run(os.Getenv("PORT"))
}
