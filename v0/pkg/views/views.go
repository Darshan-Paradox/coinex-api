package views

type Coin struct {
	Code string `json:"coin"`
	Name string `json:"name"`
}

type CoinList map[string]Coin

type Price struct {
	Coin     string  `json:"coin"`
	Currency string  `json:"currency"`
	Price    float64 `json:"price"`
}

type PriceResponseData map[string]map[string]float64

type PriceResponse struct {
	Data map[string]map[string]float64 `json:"data"`
}
