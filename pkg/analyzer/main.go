package analyzer

import (
	"os"
	"strconv"
)

type ArbMarket struct {
	Bid float64
	Ask float64
}

type MarketCalculator interface {
	CalculateSpread() float64
	CalculateConfidence() float64
}

func (a ArbMarket) CalculateSpread() float64 {
	return ((a.Ask - a.Bid) / a.Ask) * 100
}

func (a ArbMarket) CalculateConfidence() float64 {
	spread := a.CalculateSpread()
	lowSpread, _ := strconv.ParseFloat(os.Getenv("LOW_SPREAD"), 64)
	highSpread, _ := strconv.ParseFloat(os.Getenv("HIGH_SPREAD"), 64)
	return (spread - lowSpread) / (highSpread - lowSpread)
}
