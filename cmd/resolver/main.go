package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kaplanmaxe/go-carb/pkg/analyzer"
	"github.com/kaplanmaxe/go-carb/pkg/kraken"
	"github.com/kaplanmaxe/go-carb/pkg/quadriga"
)

func main() {
	godotenv.Load()
	kraken := kraken.API{
		APIKey:    os.Getenv("KRAKEN_API_KEY"),
		APISecret: os.Getenv("KRAKEN_API_SECRET"),
		Market:    "XBTCAD",
	}
	quadriga := quadriga.API{
		APIKey:    os.Getenv("QUADRIGA_API_KEY"),
		APISecret: os.Getenv("QUADRIGA_API_SECRET"),
		Market:    "btc_cad",
	}
	krakenMarket := kraken.GetMarket()
	quadrigaMarket := quadriga.GetMarket()

	bid, _ := strconv.ParseFloat(krakenMarket.Bid, 64)
	ask, _ := strconv.ParseFloat(quadrigaMarket.Ask, 64)
	arbMarket := analyzer.ArbMarket{
		Bid: bid,
		Ask: ask,
	}

	confidence := arbMarket.CalculateConfidence()

	if confidence <= 0 {
		fmt.Println("Spread too low")
		return
	}

	quadrigaBalance := quadriga.GetBalance()

	quadrigaBTCFloat, _ := strconv.ParseFloat(quadrigaBalance.BTC, 64)
	tradeAmount := fmt.Sprintf("%.3f", quadrigaBTCFloat*confidence)

	krakenTrade, err := kraken.MarketBuy(tradeAmount)
	if err != nil {
		log.Fatal("error occurred on kraken trade", err, tradeAmount)
	}
	quadrigaTrade := quadriga.MarketSell(tradeAmount)

	if quadrigaTrade.Amount != "" && len(krakenTrade.Txid) > 0 {
		fmt.Println("Arb successful!")
	}
}
