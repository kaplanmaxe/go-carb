package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kaplanmaxe/go-carb/pkg/resolver"
)

// API is a struct that queries the Kraken API
type API struct {
	APIKey    string
	APISecret string
	Market    string
}

type krakenResponseGeneric struct {
	Errors []string    `json:"error"`
	Result interface{} `json:"result"`
}

type krakenResponse struct {
	Errors []string             `json:"error"`
	Result krakenMarketResponse `json:"result"`
}

type krakenMarketResponse struct {
	Market krakenMarket `json:"XXBTZCAD"` // figure out how to make this dynamic
}

type krakenMarket struct {
	Ask  []string `json:"a"`
	Bid  []string `json:"b"`
	Last []string `json:"c"`
}

type krakenOrder struct {
	Pair      string `json:"pair"`
	OrderType string `json:"ordertype"`
	Price     string `json:"price"`
	Nonce     string `json:"nonce"`
}

// GetMarket returns a market from Kraken
func (a API) GetMarket() resolver.Market {
	res, err := http.Get("https://api.kraken.com/0/public/Ticker?pair=" + a.Market)
	var response krakenResponse
	if err != nil {
		panic("err")
	}
	json.NewDecoder(res.Body).Decode(&response)
	return resolver.Market{
		Ask:  response.Result.Market.Ask[0],
		Bid:  response.Result.Market.Bid[0],
		Last: response.Result.Market.Last[0],
	}
}

func (a API) MarketBuy(amount string) krakenResponseGeneric {
	client := &http.Client{}
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())

	order := url.Values{
		"pair":      {"XBTUSD"},
		"ordertype": {"market"},
		// "price":     {amount},
		"type":   {"buy"},
		"volume": {"0.002"},
		"nonce":  {nonce},
	}
	req, _ := http.NewRequest("POST", "https://api.kraken.com/0/private/AddOrder", strings.NewReader(order.Encode()))
	sig := a.generateSignature("/0/private/AddOrder", nonce, order)
	fmt.Println(a.APIKey)
	req.Header.Add("API-Key", a.APIKey)
	req.Header.Add("Api-Sign", sig)
	req.Header.Add("User-Agent", "go-carb")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	var response krakenResponseGeneric
	json.Unmarshal(body, &response)
	return response

}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func (a API) generateSignature(uri string, nonce string, order url.Values) string {
	shaSum := getSha256([]byte(nonce + order.Encode()))
	secret, _ := base64.StdEncoding.DecodeString(a.APISecret)
	macSum := getHMacSha512(append([]byte(uri), shaSum...), []byte(secret))
	return base64.StdEncoding.EncodeToString(macSum)
}
