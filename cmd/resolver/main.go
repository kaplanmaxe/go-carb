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

	resp := kraken.MarketBuy("0.001")
	fmt.Println(resp)
	fmt.Println(krakenMarket, quadrigaMarket, arbMarket.CalculateSpread())

	// resp := resolver.GetMarket("USDT-DASH")
	// market := analyzer.Market{Bid: resp.Result.Bid, Ask: resp.Result.Ask}
	// bittrex.LimitSell(market.Ask)
	// fmt.Println(spread)
}
