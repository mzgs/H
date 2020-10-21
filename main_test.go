package H

import "testing"

func TestUrl(t *testing.T) {

	P(UrlString("deneme asca Ğ Iİ c "))

}

func TestArrayFromUrl(t *testing.T) {

	var m []BtcturkSymbol

	err := ArrayFromUrl("https://api.btcturk.com/api/v2/ticker", "data", &m)

	if err != nil {
		t.Error(err)
	}
}

func TestPL(t *testing.T) {

	PL("mustafa", "dene")

}

type BitexenSymbol struct {
	Market struct {
		MarketCode          string `json:"market_code"`
		BaseCurrencyCode    string `json:"base_currency_code"`
		CounterCurrencyCode string `json:"counter_currency_code"`
	} `json:"market"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	LastPrice string `json:"last_price"`
	LastSize  string `json:"last_size"`
	Volume24H string `json:"volume_24h"`
	Change24H string `json:"change_24h"`
	Low24H    string `json:"low_24h"`
	High24H   string `json:"high_24h"`
	Avg24H    string `json:"avg_24h"`
	Timestamp string `json:"timestamp"`
}

type BtcturkSymbol struct {
	Pair              string  `json:"pair"`
	PairNormalized    string  `json:"pairNormalized"`
	Timestamp         int64   `json:"timestamp"`
	Last              float64 `json:"last"`
	High              float64 `json:"high"`
	Low               float64 `json:"low"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Open              float64 `json:"open"`
	Volume            float64 `json:"volume"`
	Average           float64 `json:"average"`
	Daily             float64 `json:"daily"`
	DailyPercent      float64 `json:"dailyPercent"`
	DenominatorSymbol string  `json:"denominatorSymbol"`
	NumeratorSymbol   string  `json:"numeratorSymbol"`
	Order             int     `json:"order"`
}
