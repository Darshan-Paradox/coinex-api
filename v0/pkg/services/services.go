package services

import (
	"fmt"

	"coinex-api/v0/internal/coinbase"
	"coinex-api/v0/pkg/views"
)

func GetAllCoins() (views.CoinList, error) {
	coinList, err := coinbase.FetchCoinsList()
	if err != nil {
		fmt.Println("Error in getting all coin details", err)
		return nil, err
	}
	return coinList, nil
}

func GetPrice(coinCode string, currency string) (views.Price, error) {
	coinList, err := coinbase.FetchCoinsList()
	if err != nil {
		fmt.Println("Error in getting coin details", err)
		return views.Price{}, err
	}
	price, err := coinbase.FetchCoinPrice(coinList[coinCode], currency)
	if err != nil {
		fmt.Println("Error in getting price or wrong coin code", err)
		return views.Price{}, err
	}
	return views.Price{Coin: coinCode, Currency: currency, Price: price}, nil
}
