package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"coinex-api/v0/db"
	"coinex-api/v0/pkg/services"
)

func GetAllCoins(c *gin.Context) {
	coinList, err := services.GetAllCoins()
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coinList)
}

func GetCoin(c *gin.Context) {
	coinCode := strings.ToUpper(c.Param("coin"))

	expiry, _ := db.Repository.GetExpiry(coinCode)
	if time.Now().Before(expiry) {
		coin, err := db.Repository.GetCoin(coinCode)
		if err == nil {
			c.JSON(http.StatusOK, coin)
			return
		}
	}

	coinList, err := services.GetAllCoins()
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if coin, ok := coinList[coinCode]; ok {
		db.Repository.SetCoin(coin)
		c.JSON(http.StatusOK, coinList[coinCode])
		return
	}

	fmt.Println("ERROR:: Coin doesn't exist...")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Coin doesn't exists"})
}

func GetCoinPrice(c *gin.Context) {
	coinCode, currency := strings.ToUpper(c.Param("coin")), strings.ToUpper(c.Param("currency"))

	expiry, _ := db.Repository.GetExpiry(coinCode)
	if time.Now().Before(expiry) {
		price, err := db.Repository.GetPrice(coinCode, currency)
		if err == nil {
			c.JSON(http.StatusOK, price)
			return
		}
	}

	price, err := services.GetPrice(coinCode, currency)
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Repository.GetCoin(price.Coin)

	if err != nil {
		coinList, err := services.GetAllCoins()
		if err != nil {
			fmt.Println("ERROR:: Unable to fetch data...", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		db.Repository.SetCoin(coinList[coinCode])
	}

	db.Repository.SetPrice(price)
	c.JSON(http.StatusOK, price)
}
