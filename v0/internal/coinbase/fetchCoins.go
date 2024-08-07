package coinbase

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"coinex-api/v0/pkg/views"
)

func FetchCoinsList() (views.CoinList, error) {
	response, err := http.Get(os.Getenv("COINBASE_URL") + "/currencies/crypto")
	if err != nil {
		log.Println("Internal Error: couldn't retrieving data...", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Internal Error: corrupted response", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Internal Error: corrupted data", err)
		return nil, err
	}

	allCoinsList := make(views.CoinList)

	if list, ok := data["data"].([]interface{}); ok {
		for _, val := range list {
			if coin, ok := val.(map[string]interface{}); ok {
				allCoinsList[coin["code"].(string)] = views.Coin{
					Code: coin["code"].(string),
					Name: coin["name"].(string),
				}
			}
		}
		return allCoinsList, nil
	}
	return nil, &InvalidStructure{}
}
