package cookies

import (
	"os"
	"fmt"
	"time"
	"errors"
	"strings"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"

	"coinex-api/v0/pkg/views"
)

var Repository Cookies

type Cookies struct {
    ctx *gin.Context
}

func (cookies *Cookies) Init(ctx *gin.Context) {
	cookies.ctx = ctx
}

func (cookies Cookies) GetAllCoin() (views.CoinList, error) {
	coins := make(views.CoinList)

	cookie, err := cookies.ctx.Cookie("coins")
	if err != nil {
		fmt.Println("Error: cookie not found", err)
		return coins, err
	}

	var coinList map[string]interface{}
	err = json.Unmarshal([]byte(cookie), &coins)

	if err != nil {
		fmt.Println("Error: corrupt data in cookies", err)
		return coins, err
	}

	for _, code := range coinList["data"].([]string) {
		cookie, err = cookies.ctx.Cookie(code)
		var coin map[string]string
		err = json.Unmarshal([]byte(cookie), &coin)
		if err != nil {
			fmt.Println("Error: corrupt data in cookies", err)
			return coins, err
		}
		coins[code] = views.Coin{Code: code, Name: coin["name"]}
	}

	return coins, err
}

func (cookies Cookies) GetCoin(coinCode string) (views.Coin, error) {

	coinKey := fmt.Sprintf("%s-%s", coinCode, os.Getenv("CURRENCY"))
	cookie, err := cookies.ctx.Cookie(coinKey+"-name")
	if err != nil {
		fmt.Println("Error: cookie not found", err)
		return views.Coin{}, err
	}

	return views.Coin{Code: coinCode, Name: cookie}, err
}

func (cookies Cookies) GetCoinInCurrency(coinCode, currency string) (views.Coin, error) {

	coinKey := fmt.Sprintf("%s-%s", coinCode, currency)
	cookie, err := cookies.ctx.Cookie(coinKey+"-name")
	if err != nil {
		fmt.Println("Error: cookie not found", err)
		return views.Coin{}, err
	}

	return views.Coin{Code: coinCode, Name: cookie}, err
}

func (cookies Cookies) GetPriceIn(coinCode, currency string) (views.Price, error) {

	coinKey := fmt.Sprintf("%s-%s", coinCode, currency)
	cookie, err := cookies.ctx.Cookie(coinKey+"-price")
	if err != nil {
		fmt.Println("Error: cookie not found", err)
		return views.Price{}, err
	}

	price, err := strconv.ParseFloat(cookie, 32)
	if err != nil {
		fmt.Println("Wrong format of coin price", err)
		return views.Price{}, err
	}
	if price < 0 {
		fmt.Println("Price doesn't exists in database, fetching from API...")
		return views.Price{}, errors.New("Price doesn't exists")
	}

	return views.Price{Coin: coinCode, Currency: currency, Price: price}, err
}

func (cookies Cookies) GetPrice(coinCode string) (views.PriceResponse, error) {
	currencies := strings.Split(os.Getenv("CURRENCY_LIST"), "-")
	response := make(map[string]float64)

	for _, currency := range currencies {
		coinKey := fmt.Sprintf("%s-%s", coinCode, currency)
		cookie, err := cookies.ctx.Cookie(coinKey+"-price")
		if err != nil {
			fmt.Println("Error: cookie not found", err)
			return views.PriceResponse{}, err
		}

		price, err := strconv.ParseFloat(cookie, 32)
		if err != nil {
			fmt.Println("Wrong format of coin price", err)
			return views.PriceResponse{}, err
		}
		if price < 0 {
			fmt.Println("Price doesn't exists in database, fetching from API...")
			return views.PriceResponse{}, errors.New("Price doesn't exists")
		}
		response[currency] = price
	}

	return views.PriceResponse{Data: views.PriceResponseData{coinCode: response}}, nil
}

func (cookies Cookies) SetCoin(coin views.Coin) error {
	EXPIRY, err := time.ParseDuration(os.Getenv("EXPIRY"))
	if err != nil {
		fmt.Println("Invalid time format in environment variable 'EXPIRY'", err)
		return err
	}

	currencies := strings.Split(os.Getenv("CURRENCY_LIST"), "-")
	for _, currency := range currencies {
		coinKey := fmt.Sprintf("%s-%s", coin.Code, currency)
		cookies.ctx.SetCookie(coinKey+"-name", coin.Name, int(EXPIRY.Milliseconds()), "/", "localhost", false, true)
		cookies.ctx.SetCookie(coinKey+"-price", "-1.0", int(EXPIRY.Milliseconds()), "/", "localhost", false, true)
		cookies.ctx.SetCookie("coins", coinKey+"-name", int(EXPIRY.Milliseconds()), "/", "localhost", false, true)
	}

	return nil 
}

func (cookies Cookies) SetPrice(price views.Price) error {
	EXPIRY, err := time.ParseDuration(os.Getenv("EXPIRY"))
	if err != nil {
		fmt.Println("Invalid time format in environment variable 'EXPIRY'", err)
		return err
	}

	coinKey := fmt.Sprintf("%s-%s", price.Coin, price.Currency)
	cookies.ctx.SetCookie(coinKey+"-price", fmt.Sprintf("%f", price.Price), int(EXPIRY.Milliseconds()), "/", "localhost", false, true)
	return nil
}

func (cookies Cookies) Close() {}
