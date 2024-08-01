package routes

import (
	"github.com/gin-gonic/gin"

	"coinex-api/v0/pkg/handlers"
	"coinex-api/v0/logger"
)

type Router struct {
	router *gin.Engine
}

var RouterInst Router

func (r *Router) Init() {
	r.router = gin.Default()

	r.router.Use(logger.LoggerInst.RequestHandler)
	r.router.Use(logger.LoggerInst.ResponseHandler)

	r.router.GET("/coins", handlers.GetAllCoins)
	r.router.GET("/price/:coin/:currency", handlers.GetCoinPrice)
	r.router.GET("/:coin", handlers.GetCoin)
}

func (r *Router) Run(port string) {
	r.router.Run(":" + port)
}
