package quadriga

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

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

// BalanceResponse is response from Balance endpoint
type BalanceResponse struct {
	BTC string `json:"btc_available"`
	CAD string `json:"cad_available"`
}

// TradeResponse is response from Buy/Sell endpoint
type TradeResponse struct {
	Amount string `json:"amount"`
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

// GetBalance returns balance from Quadriga
func (a API) GetBalance() *BalanceResponse {
	bal := a.makeRequest("/balance", make(map[string]string), &BalanceResponse{})
	return bal.(*BalanceResponse)
}

// MarketSell performs a market sell
func (a API) MarketSell(amount string) *TradeResponse {
	payload := make(map[string]string)
	payload["amount"] = amount
	payload["book"] = a.Market

	trade := a.makeRequest("/sell", payload, &TradeResponse{})
	return trade.(*TradeResponse)
}

func (a API) makeRequest(uri string, payload map[string]string, returnTyp interface{}) interface{} {
	urlPrefix := "https://api.quadrigacx.com/v2"
	payload["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())
	payload["key"] = a.APIKey
	payload["signature"] = a.generateSignature(payload["nonce"])
	client := &http.Client{}
	jsonString, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", urlPrefix+uri, strings.NewReader(string(jsonString)))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var returnData interface{}
	returnData = returnTyp
	json.Unmarshal(body, &returnData)
	return returnData
}

func (a API) generateSignature(nonce string) string {
	mac := hmac.New(sha256.New, []byte(a.APISecret))
	sigString := nonce + os.Getenv("QUADRIGA_CLIENT_ID") + a.APIKey
	mac.Write([]byte(sigString))
	return hex.EncodeToString(mac.Sum(nil))
}
