package main

import "math"

type StockTicker int

const (
	META StockTicker = iota
	APPL
)

type LiveStockService interface {
	CurrentPrice(ticker StockTicker) float64
}

// Stock represents a stock that retrieves its current value from a live stock service.
// Before performing calculations with the stock value, the currentValue field should be
// updated from the service to prevent miscalculations if the price changes.
type Stock struct {
	Ticker       StockTicker
	Quantity     int
	Service      LiveStockService
	currentValue float64 // latest price per share (private)
	totalValue   float64 // Quantity * currentValue (private)
}

// Freeze updates the stock's currentValue and totalValue fields
// using the live stock service. This should be called before performing
// calculations to ensure the values do not change during the calculation.
func (s *Stock) Freeze() {
	s.currentValue = s.Service.CurrentPrice(s.Ticker)
	s.totalValue = s.currentValue * float64(s.Quantity)
}

// GetAmountToBuyWith calculates how many shares can be bought with the given
// amount of money. The result is rounded to 4 decimal places.
func (s *Stock) GetAmountToBuyWith(diff float64) float64 {
	return math.Round(diff/s.currentValue*1e5) / 1e5
}

// NewStock creates and returns a new Stock with the specified ticker,
// quantity, and live stock service.
func NewStock(ticker StockTicker, quantity int, service LiveStockService) *Stock {
	return &Stock{Ticker: ticker, Quantity: quantity, Service: service}
}
