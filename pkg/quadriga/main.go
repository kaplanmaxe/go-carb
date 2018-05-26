package quadriga

import (
	"encoding/json"
	"net/http"

	"github.com/kaplanmaxe/go-carb/pkg/resolver"
)

// API is a struct that queries quadriga's API
type API struct {
	APIKey    string
	APISecret string
	Market    string
}

type quadrigaResponse struct {
	Ask  string `json:"ask"`
	Bid  string `json:"bid"`
	Last string `json:"last"`
}

// GetMarket returns a market from Quadriga
func (a API) GetMarket() resolver.Market {
	res, err := http.Get("https://api.quadrigacx.com/v2/ticker?book=" + a.Market)
	var response quadrigaResponse
	if err != nil {
		panic("err")
	}
	json.NewDecoder(res.Body).Decode(&response)
	return resolver.Market{
		Ask:  response.Ask,
		Bid:  response.Bid,
		Last: response.Last,
	}
}
