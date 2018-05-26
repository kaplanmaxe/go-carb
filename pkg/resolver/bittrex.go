package resolver

import (
	"encoding/json"
	"net/http"
)

// BittrexResp is result struct for Bittrex API
type BittrexResp struct {
	Result bittrexResult `json:"result"`
}

type bittrexResult struct {
	Bid  float64 `json:"Bid"`
	Ask  float64 `json:"Ask"`
	Last float64 `json:"Last"`
}

// GetMarket returns best bid, ask, and last price
func GetMarket(market string) BittrexResp {
	var response BittrexResp
	res, _ := http.Get("https://bittrex.com/api/v1.1/public/getticker?market=" + market)
	json.NewDecoder(res.Body).Decode(&response)
	return response
}
