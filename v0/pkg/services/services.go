package services

import (
	"fmt"

	"coinex-api/v0/internal/crypto"
	"coinex-api/v0/pkg/views"
)

func GetAllCoins() (views.CoinList, error) {
	coinList, err := crypto.FetchCoinsList()
	if err != nil {
		fmt.Println("Error in getting all coin details", err)
		return nil, err
	}
	return coinList, nil
}

func GetPrice(coinCode string, currency string) (views.Price, error) {
	coinList, err := crypto.FetchCoinsList()
	if err != nil {
		fmt.Println("Error in getting coin details", err)
		return views.Price{}, err
	}
	price, err := crypto.FetchCoinPrice(coinList[coinCode], currency)
	if err != nil {
		fmt.Println("Error in getting price or wrong coin code", err)
		return views.Price{}, err
	}
	return views.Price{Coin: coinCode, Currency: currency, Price: price}, nil
}
