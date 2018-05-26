package analyzer

const lowSpread = 0.75
const highSpread = 1.25

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
	return (spread - lowSpread) / (highSpread - lowSpread)
}
