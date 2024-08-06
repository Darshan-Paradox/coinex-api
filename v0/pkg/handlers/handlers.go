package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"coinex-api/v0/cache"
	"coinex-api/v0/cookies"
	"coinex-api/v0/pkg/services"
	"coinex-api/v0/pkg/views"
)

func GetAllCoins(c *gin.Context) {
	cookies.Repository.Init(c)
	coinList, err := services.GetAllCoins()
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coinList)
}

func GetCoin(c *gin.Context) {
	cookies.Repository.Init(c)
	coinCode := strings.ToUpper(c.Param("coin"))

	coin, err := cache.Repository.Store.GetCoin(coinCode)
	if err == nil {
		c.JSON(http.StatusOK, coin)
		return
	}

	coinList, err := services.GetAllCoins()
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if coin, ok := coinList[coinCode]; ok {
		cache.Repository.Store.SetCoin(coin)
		c.JSON(http.StatusOK, coinList[coinCode])
		return
	}

	fmt.Println("ERROR:: Coin doesn't exist...")
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Coin doesn't exists"})
}

func GetCoinPrice(c *gin.Context) {
	cookies.Repository.Init(c)
	coinCode, currency := strings.ToUpper(c.Param("coin")), strings.ToUpper(c.Param("currency"))

	price, err := cache.Repository.Store.GetPriceIn(coinCode, currency)
	if err == nil {
		c.JSON(http.StatusOK, price)
		return
	}

	price, err = services.GetPrice(coinCode, currency)
	if err != nil {
		fmt.Println("ERROR:: Unable to fetch data...", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = cache.Repository.Store.GetCoin(price.Coin)

	if err != nil {
		coinList, err := services.GetAllCoins()
		if err != nil {
			fmt.Println("ERROR:: Unable to fetch data...", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		cache.Repository.Store.SetCoin(coinList[coinCode])
	}

	cache.Repository.Store.SetPrice(price)
	c.JSON(http.StatusOK, price)
}

func GetPrice(c *gin.Context) {
	cookies.Repository.Init(c)
	coinCode := strings.ToUpper(c.Param("coin"))

	price, err := cache.Repository.Store.GetPrice(coinCode)
	if err == nil {
		c.JSON(http.StatusOK, price)
		return
	}

	currencies := strings.Split(os.Getenv("CURRENCY_LIST"), "-")

	response := make(map[string]float64)
	for _, currency := range currencies {
		price, err := services.GetPrice(coinCode, currency)
		if err != nil {
			fmt.Println("ERROR:: Unable to fetch data...", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = cache.Repository.Store.GetCoinInCurrency(price.Coin, currency)

		if err != nil {
			coinList, err := services.GetAllCoins()
			if err != nil {
				fmt.Println("ERROR:: Unable to fetch data...", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			cache.Repository.Store.SetCoin(coinList[coinCode])
		}
		cache.Repository.Store.SetPrice(price)

		response[currency] = price.Price
	}

	c.JSON(http.StatusOK, views.PriceResponse{Data: views.PriceResponseData{coinCode: response}})
}
