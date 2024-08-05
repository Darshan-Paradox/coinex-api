package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"coinex-api/v0/pkg/views"
)

func FetchCoinPrice(coin views.Coin, currency string) (float64, error) {
	//get the url from env file
	queryLink := fmt.Sprintf(os.Getenv("BASE_URL")+"/prices/%v-%v/spot", coin.Code, currency)
	response, err := http.Get(queryLink)
	if err != nil {
		fmt.Println("Internal Error: couldn't retrieving data...", err)
		return -1.0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Internal Error: corrupted response", err)
		return -1.0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Internal Error: corrupted data", err)
		return -1.0, err
	}

	if val, ok := data["data"].(interface{}); ok {
		if price, ok := val.(map[string]interface{}); ok {
			amount, _ := strconv.ParseFloat(price["amount"].(string), 32)
			return amount, nil
		}
	}

	return -1.0, &InvalidStructure{}
}
