package main

import (
	"fmt"
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
		APIKey:    os.Getenv("KRAKEN_API_KEY"),
		APISecret: os.Getenv("KRAKEN_API_SECRET"),
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

	fmt.Println(arbMarket.CalculateSpread(), arbMarket.CalculateConfidence())

	resp, _ := kraken.GetTradeBalance()
	fmt.Println(resp.ZCAD)
}
