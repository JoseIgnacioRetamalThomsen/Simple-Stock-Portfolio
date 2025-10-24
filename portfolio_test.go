package main

import (
	"reflect"
	"testing"
)

type MockLiveStockService struct{}

func (m MockLiveStockService) CurrentPrice(ticker StockTicker) float64 {
	switch ticker {
	case META:
		return 150
	case APPL:
		return 200
	default:
		return 0
	}
}

func TestPortfolioStruct_TotalValue(t *testing.T) {
	var service MockLiveStockService
	type fields struct {
		Stocks          map[StockTicker]*Stock
		AllocatedStocks map[StockTicker]float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{"returns_0_when_empty", fields{map[StockTicker]*Stock{}, make(map[StockTicker]float64)}, 0.0},
		{"returns_correct_total_when_one_stock", fields{
			map[StockTicker]*Stock{META: NewStock(META, 5, service)},
			make(map[StockTicker]float64)}, 150 * 5},
		{"returns_correct_total_when_two_stocks", fields{map[StockTicker]*Stock{
			META: NewStock(META, 2, service),
			APPL: NewStock(APPL, 3, service)},
			make(map[StockTicker]float64)}, 150*2 + 200*3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio := &PortfolioStruct{
				Stocks:          tt.fields.Stocks,
				AllocatedStocks: tt.fields.AllocatedStocks,
			}
			if got := portfolio.totalValue(); got != tt.want {
				t.Errorf("TotalValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPortfolioStruct_Rebalance(t *testing.T) {
	var service MockLiveStockService
	type fields struct {
		Stocks          map[StockTicker]*Stock
		AllocatedStocks map[StockTicker]float64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[StockTicker]float64
	}{
		{"returns_empty_when_empty", fields{map[StockTicker]*Stock{}, make(map[StockTicker]float64)}, make(map[StockTicker]float64)},
		{"rebalance_sells_meta_buys_appl", fields{map[StockTicker]*Stock{
			META: NewStock(META, 20, service),
			APPL: NewStock(APPL, 0, service)},
			map[StockTicker]float64{
				META: 0,
				APPL: 1,
			}}, map[StockTicker]float64{
			META: -20,
			APPL: 15,
		}},
		{"adjusts_only_nonempty_stocks", fields{map[StockTicker]*Stock{
			META: NewStock(META, 0, service),
			APPL: NewStock(APPL, 3, service)},
			map[StockTicker]float64{
				META: 0.5,
				APPL: 0.5,
			}}, map[StockTicker]float64{
			META: 2,
			APPL: -1.5,
		}},
		{"partial_rebalance_for_nonempty_stocks", fields{map[StockTicker]*Stock{
			META: NewStock(META, 2, service),
			APPL: NewStock(APPL, 4, service)},
			map[StockTicker]float64{
				META: 0.5,
				APPL: 0.5,
			}}, map[StockTicker]float64{
			META: 1.66667,
			APPL: -1.25,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			portfolio := &PortfolioStruct{
				Stocks:          tt.fields.Stocks,
				AllocatedStocks: tt.fields.AllocatedStocks,
			}
			if got := portfolio.Rebalance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rebalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
