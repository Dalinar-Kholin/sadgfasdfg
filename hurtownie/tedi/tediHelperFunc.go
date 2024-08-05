package tedi

import (
	"math/rand"
	"time"
)

func RandomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	stringChars := make([]byte, length)
	for i := range stringChars {
		stringChars[i] = chars[rand.Intn(len(chars))]
	}
	return string(stringChars)
}

// Stock represents the stock structure in the JSON.
type Stock struct {
	QuantityAvailable float64 `json:"quantity_available"`
	StockLocationName string  `json:"stock_location_name"`
	QuantityAllocated string  `json:"quantity_allocated"`
}

// Result represents a single product result in the JSON.
type Result struct {
	CumulativeUnitRatioSplitter string  `json:"cumulative_unit_ratio_splitter"`
	LogisticMinimum             string  `json:"logistic_minimum"`
	IsVisible                   bool    `json:"is_visible"`
	FinalPrice                  float64 `json:"final_price"`
	Stocks                      []Stock `json:"stocks"`
}

// APIResponse represents the root structure of the JSON.
type ProductResponse struct {
	Count   int      `json:"count"`
	Results []Result `json:"results"`
}
