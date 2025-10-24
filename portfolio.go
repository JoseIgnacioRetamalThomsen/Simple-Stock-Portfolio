// Package main implements a simple stock portfolio with rebalancing functionality.
package main

type Portfolio interface {
	Rebalance() map[string]int
}

// PortfolioStruct represents a portfolio containing multiple stocks
// and their target allocation percentages.
type PortfolioStruct struct {
	Stocks          map[StockTicker]*Stock  // current stocks in the portfolio
	AllocatedStocks map[StockTicker]float64 // target allocation per stock (e.g. 0.4 = 40%)
}

// Freeze all stocks to perform calculations without inconsistencies
func (portfolio *PortfolioStruct) freezeStocks() {
	for _, stock := range portfolio.Stocks {
		stock.Freeze()
	}
}

// totalValue calculates and returns the total market value of all stocks
// in the portfolio.
func (portfolio *PortfolioStruct) totalValue() float64 {
	total := 0.0
	for _, stock := range portfolio.Stocks {
		total += stock.totalValue
	}
	return total
}

// GeStockCurrentTotalValue Get current value of a stock
func (portfolio *PortfolioStruct) GeStockCurrentTotalValue(ticker StockTicker) float64 {
	return portfolio.Stocks[ticker].totalValue
}

// Rebalance returns the amounts of stocks that should be sold and bought
// to achieve the target allocation. Negative values indicate selling, positive
// indicate buying. Amounts are rounded to 4 decimal places.
// freezeStocks is called to ensure stock prices remain consistent during calculation.

func (portfolio *PortfolioStruct) Rebalance() map[StockTicker]float64 {
	portfolio.freezeStocks()
	total := portfolio.totalValue()
	adjustedStocksValues := make(map[StockTicker]float64)
	for ticker, targetPercent := range portfolio.AllocatedStocks {
		currentValue := portfolio.GeStockCurrentTotalValue(ticker)
		targetValue := total * targetPercent
		diff := targetValue - currentValue
		adjustedStocksValues[ticker] = portfolio.Stocks[ticker].GetAmountToBuyWith(diff)
	}
	return adjustedStocksValues
}
