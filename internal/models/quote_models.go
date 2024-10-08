// Package models contains the models for the Moneybots API
package models

// QuoteResponse is the response for the quote API
type QuoteResponse struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// OHLC is the OHLC data for a given instrument
type OHLC struct {
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// Depth is the depth data for a given instrument
type Depth struct {
	Buy  [5]DepthItem `json:"buy"`
	Sell [5]DepthItem `json:"sell"`
}

// DepthItem is the depth item for a given instrument
type DepthItem struct {
	Price    float64 `json:"price"`
	Quantity uint32  `json:"quantity"`
	Orders   uint32  `json:"orders"`
}

// QuoteData is the quote data for a given instrument
type QuoteData struct {
	Instrument         string  `json:"instrument"`
	Mode               string  `json:"mode"`
	InstrumentToken    uint32  `json:"instrument_token"`
	IsTradable         bool    `json:"is_tradable"`
	IsIndex            bool    `json:"is_index"`
	Timestamp          string  `json:"timestamp"`
	LastTradeTime      string  `json:"last_trade_time"`
	LastPrice          float64 `json:"last_price"`
	LastTradedQuantity uint32  `json:"last_traded_quantity"`
	TotalBuyQuantity   uint32  `json:"total_buy_quantity"`
	TotalSellQuantity  uint32  `json:"total_sell_quantity"`
	VolumeTraded       uint32  `json:"volume"`
	// TotalBuy           uint32  `json:"total_buy"`
	// TotalSell          uint32  `json:"total_sell"`
	AverageTradePrice float64 `json:"average_price"`
	OI                uint32  `json:"oi"`
	OIDayHigh         uint32  `json:"oi_day_high"`
	OIDayLow          uint32  `json:"oi_day_low"`
	NetChange         float64 `json:"net_change"`
	OHLC              OHLC    `json:"ohlc"`
	Depth             Depth   `json:"depth"`
	UpdatedAt         string  `json:"-"`
}

// OHLCData is the OHLC data for a given instrument
type OHLCData struct {
	InstrumentToken   uint32  `json:"-"`
	LastPrice         float64 `json:"last_price"`
	VolumeTraded      uint32  `json:"volume"`
	AverageTradePrice float64 `json:"average_price"`
	Timestamp         string  `json:"timestamp"`
	LastTradeTime     string  `json:"last_trade_time"`
	OHLC              OHLC    `json:"ohlc"`
	UpdatedAt         string  `json:"-"`
}

// LTPData is the LTP data for a given instrument
type LTPData struct {
	InstrumentToken uint32  `json:"-"`
	LastPrice       float64 `json:"last_price"`
	Timestamp       string  `json:"timestamp"`
	UpdatedAt       string  `json:"-"`
}
