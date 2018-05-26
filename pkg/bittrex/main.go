package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type tradeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func formatPrice(price float64) string {
	// res := strconv.FormatFloat(float64(price), 'E', -1, 32)
	return fmt.Sprintf("%f", price-0.00005)
}

func makeRequest(uri string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", uri, nil)
	sig := apiSign(uri)
	req.Header.Add("apisign", sig)
	return client.Do(req)
}

func LimitSell(price float64) tradeResponse {
	var response tradeResponse
	fmt.Println(formatPrice(price))
	uri := "https://bittrex.com/api/v1.1/market/selllimit?apikey=" + os.Getenv("BITTREX_API_KEY") +
		"&market=USDT-DASH&quantity=0.01&rate=" + formatPrice(price) + "&nonce=" + strconv.Itoa(int(time.Now().UnixNano()))
	res, _ := makeRequest(uri)
	fmt.Println(formatPrice(price))
	json.NewDecoder(res.Body).Decode(&response)
	fmt.Println(response)
	return response
}

func apiSign(uri string) string {
	mac := hmac.New(sha512.New, []byte(os.Getenv("BITTREX_API_SECRET")))
	mac.Write([]byte(uri))
	return hex.EncodeToString(mac.Sum(nil))
}
